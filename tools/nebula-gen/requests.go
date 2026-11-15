package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func (g *generator) query(op *operation, n operationNames) error {
	if len(op.Parameters) == 0 {
		return nil
	}
	b := g.chunk("request_" + n.Group + ".go")
	fields := []modelField{}
	hasAsset, hasSubject := false, false
	for _, p := range op.Parameters {
		if p.Name == "asset_id" {
			hasAsset = true
		}
		if p.Name == "subject_id" {
			hasSubject = true
		}
	}
	if hasAsset {
		fields = append(fields, modelField{name: "Asset", typ: "*AssetRef", description: "Asset selects an asset by its identifier and class."})
	}
	if hasSubject {
		fields = append(fields, modelField{name: "Subject", typ: "*SubjectRef", description: "Subject selects a non-asset person or topic."})
	}
	for _, p := range op.Parameters {
		if p.Name == "locale" || p.Name == "asset_id" || p.Name == "subject_id" || p.Name == "asset_class" && hasAsset && !n.StandaloneClass {
			continue
		}
		s, err := g.doc.resolve(p.Schema)
		if err != nil {
			return err
		}
		typ, err := g.goType(p.Schema, op.ID+"/query/"+p.Name, n.Request+goName(p.Name), n.Group, true)
		if err != nil {
			return err
		}
		if s.Type == "object" {
			return fmt.Errorf("%s: object query parameter %s is unsupported", op.ID, p.Name)
		}
		if !p.Required {
			typ = optional(typ)
		}
		fields = append(fields, modelField{g.fieldName(n.Request, p.Name), typ, p.Name, p.Description, p.Required})
	}
	m := &model{name: n.Request, fields: fields}
	g.orderFields(m)
	fields = m.fields
	comment(b, n.Request+" carries the query parameters of "+op.ID+".")
	fmt.Fprintf(b, "type %s struct {\n", n.Request)
	for _, f := range fields {
		comment(b, f.description)
		fmt.Fprintf(b, "%s %s\n", f.name, f.typ)
	}
	fmt.Fprint(b, "}\n\n")
	fmt.Fprintf(b, "func(r *%s) query()(url.Values,error){\nv:=url.Values{}\nif r==nil{return v,validateQuery(%q,v)}\n", n.Request, op.ID)
	if hasAsset {
		fmt.Fprintf(b, "if err:=r.Asset.optionalPair(%q);err!=nil{return nil,err}\n", op.ID)
	}
	if n.StandaloneClass && hasAsset {
		subject := "nil"
		if hasSubject {
			subject = "r.Subject"
		}
		fmt.Fprintf(b, "if err:=validateClassFilter(r.Asset,r.AssetClass,%s,%q);err!=nil{return nil,err}\n", subject, op.ID)
	}
	for _, p := range op.Parameters {
		s, err := g.doc.resolve(p.Schema)
		if err != nil {
			return err
		}
		switch p.Name {
		case "locale":
			continue
		case "asset_id":
			fmt.Fprintln(b, "if r.Asset!=nil{v.Set(\"asset_id\",r.Asset.ID);v.Set(\"asset_class\",string(r.Asset.Class))}")
			continue
		case "subject_id":
			fmt.Fprintln(b, "if r.Subject!=nil{v.Set(\"subject_id\",r.Subject.ID)}")
			continue
		case "asset_class":
			if hasAsset && !n.StandaloneClass {
				continue
			}
		}
		field := g.fieldName(n.Request, p.Name)
		expr := "r." + field
		if s.Type == "array" {
			item, err := g.doc.resolve(s.Items)
			if err != nil {
				return err
			}
			if item.Type == "array" || item.Type == "object" {
				return fmt.Errorf("%s: nested array query parameter %s is unsupported", op.ID, p.Name)
			}
			if p.Explode != nil && !*p.Explode {
				fmt.Fprintf(b, "if len(%s)>0{parts:=make([]string,0,len(%s));for _,item:=range %s{parts=append(parts,%s)};v.Set(%q,strings.Join(parts,\",\"))}\n", expr, expr, expr, encode(item, "item"), p.Name)
			} else {
				fmt.Fprintf(b, "for _,item:=range %s{v.Add(%q,%s)}\n", expr, p.Name, encode(item, "item"))
			}
		} else {
			if !p.Required {
				fmt.Fprintf(b, "if %s!=nil{\n", expr)
				expr = "*" + expr
			}
			fmt.Fprintf(b, "v.Set(%q,%s)\n", p.Name, encode(s, expr))
			if !p.Required {
				fmt.Fprintln(b, "}")
			}
		}
	}
	fmt.Fprintf(b, "return v,validateQuery(%q,v)\n}\n\n", op.ID)
	return nil
}

func encode(s *schema, expr string) string {
	switch s.Type {
	case "integer":
		if s.Format == "int64" {
			return "strconv.FormatInt(" + expr + ",10)"
		}
		return "strconv.Itoa(" + expr + ")"
	case "number":
		return "strconv.FormatFloat(" + expr + ",'g',-1,64)"
	case "boolean":
		return "strconv.FormatBool(" + expr + ")"
	default:
		if s.Format == "date-time" {
			return "(" + expr + ").UTC().Format(time.RFC3339Nano)"
		}
		return "string(" + expr + ")"
	}
}

func (g *generator) queryRules() error {
	b := g.chunk("validation_rules.go")
	fmt.Fprintln(b, "var queryRules = map[string][]queryRule{")
	for _, op := range g.doc.ops {
		if op.method != "get" {
			continue
		}
		fmt.Fprintf(b, "%q:{\n", op.ID)
		for _, p := range op.Parameters {
			if p.Name == "locale" {
				continue
			}
			s, err := g.doc.resolve(p.Schema)
			if err != nil {
				return err
			}
			fmt.Fprintf(b, "{name:%q,required:%t,", p.Name, p.Required)
			if s.Type == "array" {
				if s.MinItems != nil {
					fmt.Fprintf(b, "minItems:IntPtr(%d),", *s.MinItems)
				}
				if s.MaxItems != nil {
					fmt.Fprintf(b, "maxItems:IntPtr(%d),", *s.MaxItems)
				}
				fmt.Fprintf(b, "csv:%t,", p.Explode != nil && !*p.Explode)
				s, err = g.doc.resolve(s.Items)
				if err != nil {
					return err
				}
			}
			fmt.Fprintf(b, "kind:%q,", s.Type)
			if len(s.Enum) > 0 {
				fmt.Fprint(b, "enum:[]string{")
				for _, v := range s.Enum {
					fmt.Fprintf(b, "%q,", fmt.Sprint(v))
				}
				fmt.Fprint(b, "},")
			}
			if s.Minimum != nil {
				fmt.Fprintf(b, "min:Float64Ptr(%v),", *s.Minimum)
			}
			if s.Maximum != nil {
				fmt.Fprintf(b, "max:Float64Ptr(%v),", *s.Maximum)
			}
			if s.MinLength != nil {
				fmt.Fprintf(b, "minLength:%d,", *s.MinLength)
			}
			if s.MaxLength != nil {
				fmt.Fprintf(b, "maxLength:%d,", *s.MaxLength)
			}
			fmt.Fprintln(b, "},")
		}
		fmt.Fprintln(b, "},")
	}
	fmt.Fprintln(b, "}")
	return nil
}

func (g *generator) requestBody(op *operation, n operationNames) error {
	if op.ID != "create_chat_completion" {
		return fmt.Errorf("%s: POST adapter must be configured explicitly", op.ID)
	}
	s, err := g.doc.jsonSchema(op.RequestBody)
	if err != nil {
		return err
	}
	if s.Type != "object" {
		return fmt.Errorf("%s: expected object request body", op.ID)
	}
	if _, err = g.goType(s, op.ID+"/request", n.Request, n.Group, true); err != nil {
		return err
	}
	m := g.models[n.Request]
	if s.Properties["stream"] == nil || s.Properties["stream"].Type != "boolean" {
		return fmt.Errorf("%s: expected boolean stream property", op.ID)
	}
	b := g.chunk("request_" + n.Group + ".go")
	fmt.Fprintln(b, "type chatCompletionBody struct {")
	for _, f := range m.fields {
		typ, tag := f.typ, f.jsonName
		if f.jsonName == "stream" {
			typ = "bool"
		} else if !f.required {
			tag += ",omitempty"
		}
		fmt.Fprintf(b, "%s %s `json:%q`\n", f.name, typ, tag)
	}
	fmt.Fprintln(b, "}")
	fmt.Fprintf(b, "func(r *%s) body(stream bool)chatCompletionBody{return chatCompletionBody{\n", n.Request)
	for _, f := range m.fields {
		expr := "r." + f.name
		if f.jsonName == "stream" {
			expr = "stream"
		}
		fmt.Fprintf(b, "%s:%s,\n", f.name, expr)
	}
	fmt.Fprintln(b, "}}")
	fmt.Fprintf(b, "func(r *%s) validate()error{\nif r==nil{return &ValidationError{Op:%q,Field:\"request\",Msg:\"request is required\"}}\n", n.Request, op.ID)
	if err := g.validateBody(b, s, "r", n.Request, "", op.ID, op.ID+"/request", 0); err != nil {
		return err
	}
	fmt.Fprint(b, "return nil\n}\n\n")
	return nil
}

func (g *generator) validateBody(b *bytes.Buffer, raw *schema, expr, parent, path, op, location string, depth int) error {
	if depth > 20 {
		return fmt.Errorf("%s: recursive request validation is unsupported", path)
	}
	s, err := g.doc.resolve(raw)
	if err != nil {
		return err
	}
	fail := func(condition, message string) {
		fmt.Fprintf(b, "if %s{return &ValidationError{Op:%q,Field:%q,Msg:%q}}\n", condition, op, path, message)
	}
	switch s.Type {
	case "object":
		if len(s.Properties) == 0 && len(s.AdditionalProperties) > 0 && string(s.AdditionalProperties) != "true" && string(s.AdditionalProperties) != "false" {
			var value schema
			if err := json.Unmarshal(s.AdditionalProperties, &value); err != nil {
				return err
			}
			var nested bytes.Buffer
			item := fmt.Sprintf("v%d", depth)
			if err := g.validateBody(&nested, &value, item, parent+"Entry", path+".*", op, location+"/additionalProperties", depth+1); err != nil {
				return err
			}
			if nested.Len() > 0 {
				fmt.Fprintf(b, "for _,%s:=range %s{\n%s}\n", item, expr, nested.String())
			}
		}
		for _, key := range keys(s.Properties) {
			child, err := g.doc.resolve(s.Properties[key])
			if err != nil {
				return err
			}
			f := g.fieldName(parent, key)
			childExpr := expr + "." + f
			originalExpr := childExpr
			childPath := strings.TrimPrefix(path+"."+key, ".")
			isOptional := !contains(s.Required, key) || child.Nullable
			if !isOptional && (child.Type == "object" || child.Type == "array") {
				fmt.Fprintf(b, "if %s==nil{return &ValidationError{Op:%q,Field:%q,Msg:\"field is required\"}}\n", originalExpr, op, childPath)
			}
			if isOptional && child.Type != "array" && child.Type != "object" {
				childExpr = "(*" + childExpr + ")"
			}
			childParent := g.compat.Objects[location+"/properties/"+key]
			if childParent == "" {
				childParent = parent + goName(key)
			}
			var nested bytes.Buffer
			if err := g.validateBody(&nested, child, childExpr, childParent, childPath, op, location+"/properties/"+key, depth+1); err != nil {
				return err
			}
			if nested.Len() > 0 {
				if isOptional {
					fmt.Fprintf(b, "if %s!=nil{\n", originalExpr)
				}
				b.Write(nested.Bytes())
				if isOptional {
					fmt.Fprintln(b, "}")
				}
			}
		}
	case "array":
		if s.MinItems != nil {
			fail(fmt.Sprintf("len(%s)<%d", expr, *s.MinItems), fmt.Sprintf("must contain at least %d items", *s.MinItems))
		}
		if s.MaxItems != nil {
			fail(fmt.Sprintf("len(%s)>%d", expr, *s.MaxItems), fmt.Sprintf("must contain at most %d items", *s.MaxItems))
		}
		itemType := g.compat.Objects[location+"/items"]
		if itemType == "" {
			itemType = parent + "Item"
		}
		index := fmt.Sprintf("i%d", depth)
		var nested bytes.Buffer
		if err := g.validateBody(&nested, s.Items, expr+"["+index+"]", itemType, path+"[]", op, location+"/items", depth+1); err != nil {
			return err
		}
		if nested.Len() > 0 {
			fmt.Fprintf(b, "for %s:=range %s{\n%s}\n", index, expr, nested.String())
		}
	case "string":
		if s.MinLength != nil {
			fail(fmt.Sprintf("utf8.RuneCountInString(string(%s))<%d", expr, *s.MinLength), fmt.Sprintf("minimum length is %d", *s.MinLength))
		}
		if s.MaxLength != nil {
			fail(fmt.Sprintf("utf8.RuneCountInString(string(%s))>%d", expr, *s.MaxLength), fmt.Sprintf("maximum length is %d", *s.MaxLength))
		}
	case "number", "integer":
		if s.Type == "number" {
			fail("math.IsNaN("+expr+")||math.IsInf("+expr+",0)", "must be finite")
		}
		if s.Minimum != nil {
			fail(fmt.Sprintf("%s<%v", expr, *s.Minimum), "below minimum")
		}
		if s.Maximum != nil {
			fail(fmt.Sprintf("%s>%v", expr, *s.Maximum), "above maximum")
		}
	}
	if len(s.Enum) > 0 {
		fmt.Fprint(b, "switch ")
		if s.Type == "string" {
			fmt.Fprintf(b, "string(%s)", expr)
		} else {
			fmt.Fprint(b, expr)
		}
		fmt.Fprintln(b, "{")
		for _, v := range s.Enum {
			if s.Type == "string" {
				fmt.Fprintf(b, "case %q:\n", v)
			} else {
				fmt.Fprintf(b, "case %v:\n", v)
			}
		}
		fmt.Fprintf(b, "default:return &ValidationError{Op:%q,Field:%q,Msg:\"value is not in the documented enum\"}\n}\n", op, path)
	}
	return nil
}
