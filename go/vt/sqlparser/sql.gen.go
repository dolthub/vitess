package sqlparser

func (p *parser) functionCallKeyword() (Expr, bool) {
  id, tok := p.peek()
  if id == LEFT {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: LEFT <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: LEFT openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: LEFT openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == RIGHT {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: RIGHT <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: RIGHT openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: RIGHT openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == FORMAT {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: FORMAT <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: FORMAT openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: FORMAT openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == GROUPING {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: GROUPING <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: GROUPING openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: GROUPING openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == SCHEMA {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: SCHEMA <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != ')' {
      p.fail("expected: 'function_call_keyword: SCHEMA openb <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1))}, true
  }  else if id == CONVERT {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: CONVERT <'('>', found: '" + string(tok) + "'")
    }
    var5, ok := p.convertType()
    if !ok {
      p.fail("expected: 'function_call_keyword: CONVERT openb expression ',' <convert_type>', found: 'string(tok)'")
    }
    var6, tok := p.next()
    if var6 != ')' {
      p.fail("expected: 'function_call_keyword: CONVERT openb expression ',' convert_type <')'>', found: '" + string(tok) + "'")
    }
    return &ConvertExpr{Name: string(var1), Expr: var3, Type: var5}, true
  }  else if id == CAST {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: CAST <'('>', found: '" + string(tok) + "'")
    }
    var4, tok := p.next()
    if var4 != AS {
      p.fail("expected: 'function_call_keyword: CAST openb expression <AS>', found: '" + string(tok) + "'")
    }
    var5, ok := p.convertType()
    if !ok {
      p.fail("expected: 'function_call_keyword: CAST openb expression AS <convert_type>', found: 'string(tok)'")
    }
    var6, tok := p.next()
    if var6 != ')' {
      p.fail("expected: 'function_call_keyword: CAST openb expression AS convert_type <')'>', found: '" + string(tok) + "'")
    }
    return &ConvertExpr{Name: string(var1), Expr: var3, Type: var5}, true
  }  else if id == CHAR {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: CHAR <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: CHAR openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: CHAR openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &CharExpr{Exprs: var3}, true
  }  else if id == CHAR {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: CHAR <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: CHAR openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != USING {
      p.fail("expected: 'function_call_keyword: CHAR openb argument_expression_list <USING>', found: '" + string(tok) + "'")
    }
    var6, tok := p.next()
    if var6 != ')' {
      p.fail("expected: 'function_call_keyword: CHAR openb argument_expression_list USING charset <')'>', found: '" + string(tok) + "'")
    }
    return &CharExpr{Exprs: var3, Type: var5}, true
  }  else if id == CONVERT {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: CONVERT <'('>', found: '" + string(tok) + "'")
    }
    var4, tok := p.next()
    if var4 != USING {
      p.fail("expected: 'function_call_keyword: CONVERT openb expression <USING>', found: '" + string(tok) + "'")
    }
    var6, tok := p.next()
    if var6 != ')' {
      p.fail("expected: 'function_call_keyword: CONVERT openb expression USING charset <')'>', found: '" + string(tok) + "'")
    }
    return &ConvertUsingExpr{Expr: var3, Type: var5}, true
  }  else if id == POSITION {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: POSITION <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: POSITION openb <value_expression>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != IN {
      p.fail("expected: 'function_call_keyword: POSITION openb value_expression <IN>', found: '" + string(tok) + "'")
    }
    var5, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: POSITION openb value_expression IN <value_expression>', found: 'string(tok)'")
    }
    var6, tok := p.next()
    if var6 != ')' {
      p.fail("expected: 'function_call_keyword: POSITION openb value_expression IN value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent("LOCATE"), Exprs: []SelectExpr{&AliasedExpr{Expr: var3}, &AliasedExpr{Expr: var5}}}, true
  }  else if id == INSERT {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: INSERT <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: INSERT openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: INSERT openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == SUBSTR {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: SUBSTR <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.columnName()
    if !ok {
      p.fail("expected: 'function_call_keyword: SUBSTR openb <column_name>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != FROM {
      p.fail("expected: 'function_call_keyword: SUBSTR openb column_name <FROM>', found: '" + string(tok) + "'")
    }
    var5, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: SUBSTR openb column_name FROM <value_expression>', found: 'string(tok)'")
    }
    var6, tok := p.next()
    if var6 != FOR {
      p.fail("expected: 'function_call_keyword: SUBSTR openb column_name FROM value_expression <FOR>', found: '" + string(tok) + "'")
    }
    var7, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: SUBSTR openb column_name FROM value_expression FOR <value_expression>', found: 'string(tok)'")
    }
    var8, tok := p.next()
    if var8 != ')' {
      p.fail("expected: 'function_call_keyword: SUBSTR openb column_name FROM value_expression FOR value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &SubstrExpr{Name: var3, From: var5, To: var7}, true
  }  else if id == SUBSTRING {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: SUBSTRING <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.columnName()
    if !ok {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb <column_name>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != FROM {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb column_name <FROM>', found: '" + string(tok) + "'")
    }
    var5, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb column_name FROM <value_expression>', found: 'string(tok)'")
    }
    var6, tok := p.next()
    if var6 != FOR {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb column_name FROM value_expression <FOR>', found: '" + string(tok) + "'")
    }
    var7, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb column_name FROM value_expression FOR <value_expression>', found: 'string(tok)'")
    }
    var8, tok := p.next()
    if var8 != ')' {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb column_name FROM value_expression FOR value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &SubstrExpr{Name: var3, From: var5, To: var7}, true
  }  else if id == SUBSTR {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: SUBSTR <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != STRING {
      p.fail("expected: 'function_call_keyword: SUBSTR openb <STRING>', found: '" + string(tok) + "'")
    }
    var4, tok := p.next()
    if var4 != FROM {
      p.fail("expected: 'function_call_keyword: SUBSTR openb STRING <FROM>', found: '" + string(tok) + "'")
    }
    var5, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: SUBSTR openb STRING FROM <value_expression>', found: 'string(tok)'")
    }
    var6, tok := p.next()
    if var6 != FOR {
      p.fail("expected: 'function_call_keyword: SUBSTR openb STRING FROM value_expression <FOR>', found: '" + string(tok) + "'")
    }
    var7, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: SUBSTR openb STRING FROM value_expression FOR <value_expression>', found: 'string(tok)'")
    }
    var8, tok := p.next()
    if var8 != ')' {
      p.fail("expected: 'function_call_keyword: SUBSTR openb STRING FROM value_expression FOR value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &SubstrExpr{StrVal: NewStrVal(var3), From: var5, To: var7}, true
  }  else if id == SUBSTRING {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: SUBSTRING <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != STRING {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb <STRING>', found: '" + string(tok) + "'")
    }
    var4, tok := p.next()
    if var4 != FROM {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb STRING <FROM>', found: '" + string(tok) + "'")
    }
    var5, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb STRING FROM <value_expression>', found: 'string(tok)'")
    }
    var6, tok := p.next()
    if var6 != FOR {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb STRING FROM value_expression <FOR>', found: '" + string(tok) + "'")
    }
    var7, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb STRING FROM value_expression FOR <value_expression>', found: 'string(tok)'")
    }
    var8, tok := p.next()
    if var8 != ')' {
      p.fail("expected: 'function_call_keyword: SUBSTRING openb STRING FROM value_expression FOR value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &SubstrExpr{StrVal: NewStrVal(var3), From: var5, To: var7}, true
  }  else if id == TRIM {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: TRIM <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: TRIM openb <value_expression>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: TRIM openb value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &TrimExpr{Pattern: NewStrVal([]byte(" ")), Str: var3, Dir: Both}, true
  }  else if id == TRIM {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: TRIM <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: TRIM openb <value_expression>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != FROM {
      p.fail("expected: 'function_call_keyword: TRIM openb value_expression <FROM>', found: '" + string(tok) + "'")
    }
    var5, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: TRIM openb value_expression FROM <value_expression>', found: 'string(tok)'")
    }
    var6, tok := p.next()
    if var6 != ')' {
      p.fail("expected: 'function_call_keyword: TRIM openb value_expression FROM value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &TrimExpr{Pattern: var3, Str: var5, Dir: Both}, true
  }  else if id == TRIM {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: TRIM <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != LEADING {
      p.fail("expected: 'function_call_keyword: TRIM openb <LEADING>', found: '" + string(tok) + "'")
    }
    var4, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: TRIM openb LEADING <value_expression>', found: 'string(tok)'")
    }
    var5, tok := p.next()
    if var5 != FROM {
      p.fail("expected: 'function_call_keyword: TRIM openb LEADING value_expression <FROM>', found: '" + string(tok) + "'")
    }
    var6, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: TRIM openb LEADING value_expression FROM <value_expression>', found: 'string(tok)'")
    }
    var7, tok := p.next()
    if var7 != ')' {
      p.fail("expected: 'function_call_keyword: TRIM openb LEADING value_expression FROM value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &TrimExpr{Pattern: var4, Str: var6, Dir: Leading}, true
  }  else if id == TRIM {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: TRIM <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != TRAILING {
      p.fail("expected: 'function_call_keyword: TRIM openb <TRAILING>', found: '" + string(tok) + "'")
    }
    var4, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: TRIM openb TRAILING <value_expression>', found: 'string(tok)'")
    }
    var5, tok := p.next()
    if var5 != FROM {
      p.fail("expected: 'function_call_keyword: TRIM openb TRAILING value_expression <FROM>', found: '" + string(tok) + "'")
    }
    var6, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: TRIM openb TRAILING value_expression FROM <value_expression>', found: 'string(tok)'")
    }
    var7, tok := p.next()
    if var7 != ')' {
      p.fail("expected: 'function_call_keyword: TRIM openb TRAILING value_expression FROM value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &TrimExpr{Pattern: var4, Str: var6, Dir: Trailing}, true
  }  else if id == TRIM {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: TRIM <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != BOTH {
      p.fail("expected: 'function_call_keyword: TRIM openb <BOTH>', found: '" + string(tok) + "'")
    }
    var4, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: TRIM openb BOTH <value_expression>', found: 'string(tok)'")
    }
    var5, tok := p.next()
    if var5 != FROM {
      p.fail("expected: 'function_call_keyword: TRIM openb BOTH value_expression <FROM>', found: '" + string(tok) + "'")
    }
    var6, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: TRIM openb BOTH value_expression FROM <value_expression>', found: 'string(tok)'")
    }
    var7, tok := p.next()
    if var7 != ')' {
      p.fail("expected: 'function_call_keyword: TRIM openb BOTH value_expression FROM value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &TrimExpr{Pattern: var4, Str: var6, Dir: Both}, true
  }  else if id == MATCH {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: MATCH <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: MATCH openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: MATCH openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, tok := p.next()
    if var5 != AGAINST {
      p.fail("expected: 'function_call_keyword: MATCH openb argument_expression_list closeb <AGAINST>', found: '" + string(tok) + "'")
    }
    var6, tok := p.next()
    if var6 != '(' {
      p.fail("expected: 'function_call_keyword: MATCH openb argument_expression_list closeb AGAINST <'('>', found: '" + string(tok) + "'")
    }
    var7, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_keyword: MATCH openb argument_expression_list closeb AGAINST openb <value_expression>', found: 'string(tok)'")
    }
    var8, ok := p.matchOption()
    if !ok {
      p.fail("expected: 'function_call_keyword: MATCH openb argument_expression_list closeb AGAINST openb value_expression <match_option>', found: 'string(tok)'")
    }
    var9, tok := p.next()
    if var9 != ')' {
      p.fail("expected: 'function_call_keyword: MATCH openb argument_expression_list closeb AGAINST openb value_expression match_option <')'>', found: '" + string(tok) + "'")
    }
    return &MatchExpr{Columns: var3, Expr: var7, Option: var8}, true
  }  else if id == FIRST {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: FIRST <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: FIRST openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: FIRST openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == GROUP_CONCAT {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: GROUP_CONCAT <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.distinctOpt()
    if !ok {
      p.fail("expected: 'function_call_keyword: GROUP_CONCAT openb <distinct_opt>', found: 'string(tok)'")
    }
    var4, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: GROUP_CONCAT openb distinct_opt <argument_expression_list>', found: 'string(tok)'")
    }
    var5, ok := p.orderByOpt()
    if !ok {
      p.fail("expected: 'function_call_keyword: GROUP_CONCAT openb distinct_opt argument_expression_list <order_by_opt>', found: 'string(tok)'")
    }
    var6, ok := p.separatorOpt()
    if !ok {
      p.fail("expected: 'function_call_keyword: GROUP_CONCAT openb distinct_opt argument_expression_list order_by_opt <separator_opt>', found: 'string(tok)'")
    }
    var7, tok := p.next()
    if var7 != ')' {
      p.fail("expected: 'function_call_keyword: GROUP_CONCAT openb distinct_opt argument_expression_list order_by_opt separator_opt <')'>', found: '" + string(tok) + "'")
    }
    return &GroupConcatExpr{Distinct: var3, Exprs: var4, OrderBy: var5, Separator: var6}, true
  }  else if id == CASE {
    var1, _ := p.next()
    var2, ok := p.expressionOpt()
    if !ok {
      p.fail("expected: 'function_call_keyword: CASE <expression_opt>', found: 'string(tok)'")
    }
    var3, ok := p.whenExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: CASE expression_opt <when_expression_list>', found: 'string(tok)'")
    }
    var4, ok := p.elseExpressionOpt()
    if !ok {
      p.fail("expected: 'function_call_keyword: CASE expression_opt when_expression_list <else_expression_opt>', found: 'string(tok)'")
    }
    var5, tok := p.next()
    if var5 != END {
      p.fail("expected: 'function_call_keyword: CASE expression_opt when_expression_list else_expression_opt <END>', found: '" + string(tok) + "'")
    }
    return &CaseExpr{Expr: var2, Whens: var3, Else: var4}, true
  }  else if id == VALUES {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: VALUES <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.columnName()
    if !ok {
      p.fail("expected: 'function_call_keyword: VALUES openb <column_name>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: VALUES openb column_name <')'>', found: '" + string(tok) + "'")
    }
    return &ValuesFuncExpr{Name: var3}, true
  }  else if id == VALUES {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: VALUES <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.columnNameSafeKeyword()
    if !ok {
      p.fail("expected: 'function_call_keyword: VALUES openb <column_name_safe_keyword>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: VALUES openb column_name_safe_keyword <')'>', found: '" + string(tok) + "'")
    }
    return &ValuesFuncExpr{Name: NewColName(string(var3))}, true
  }  else if id == VALUES {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: VALUES <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.nonReservedKeyword3()
    if !ok {
      p.fail("expected: 'function_call_keyword: VALUES openb <non_reserved_keyword3>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: VALUES openb non_reserved_keyword3 <')'>', found: '" + string(tok) + "'")
    }
    return &ValuesFuncExpr{Name: NewColName(string(var3))}, true
  }  else if id == REPEAT {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_keyword: REPEAT <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_keyword: REPEAT openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_keyword: REPEAT openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }
  return nil, false
}
func (p *parser) functionCallWindow() (Expr, bool) {
  id, tok := p.peek()
  if id == CUME_DIST {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: CUME_DIST <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != ')' {
      p.fail("expected: 'function_call_window: CUME_DIST openb <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Over: var4}, true
  }  else if id == DENSE_RANK {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: DENSE_RANK <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != ')' {
      p.fail("expected: 'function_call_window: DENSE_RANK openb <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Over: var4}, true
  }  else if id == FIRST_VALUE {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: FIRST_VALUE <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpression()
    if !ok {
      p.fail("expected: 'function_call_window: FIRST_VALUE openb <argument_expression>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_window: FIRST_VALUE openb argument_expression <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: SelectExprs{var3}, Over: var5}, true
  }  else if id == LAG {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: LAG <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_window: LAG openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_window: LAG openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == LAST_VALUE {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: LAST_VALUE <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpression()
    if !ok {
      p.fail("expected: 'function_call_window: LAST_VALUE openb <argument_expression>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_window: LAST_VALUE openb argument_expression <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: SelectExprs{var3}, Over: var5}, true
  }  else if id == LEAD {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: LEAD <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_window: LEAD openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_window: LEAD openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == NTH_VALUE {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: NTH_VALUE <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_window: NTH_VALUE openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_window: NTH_VALUE openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == NTILE {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: NTILE <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != ')' {
      p.fail("expected: 'function_call_window: NTILE openb <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Over: var4}, true
  }  else if id == PERCENT_RANK {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: PERCENT_RANK <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != ')' {
      p.fail("expected: 'function_call_window: PERCENT_RANK openb <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Over: var4}, true
  }  else if id == RANK {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: RANK <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != ')' {
      p.fail("expected: 'function_call_window: RANK openb <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Over: var4}, true
  }  else if id == ROW_NUMBER {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_window: ROW_NUMBER <'('>', found: '" + string(tok) + "'")
    }
    var3, tok := p.next()
    if var3 != ')' {
      p.fail("expected: 'function_call_window: ROW_NUMBER openb <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Over: var4}, true
  }
  return nil, false
}
func (p *parser) functionCallAggregateWithWindow() (Expr, bool) {
  id, tok := p.peek()
  if id == MAX {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: MAX <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.distinctOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: MAX openb <distinct_opt>', found: 'string(tok)'")
    }
    var4, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: MAX openb distinct_opt <argument_expression_list>', found: 'string(tok)'")
    }
    var5, tok := p.next()
    if var5 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: MAX openb distinct_opt argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var6, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: MAX openb distinct_opt argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var4, Distinct: var3 == DistinctStr, Over: var6}, true
  }  else if id == AVG {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: AVG <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.distinctOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: AVG openb <distinct_opt>', found: 'string(tok)'")
    }
    var4, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: AVG openb distinct_opt <argument_expression_list>', found: 'string(tok)'")
    }
    var5, tok := p.next()
    if var5 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: AVG openb distinct_opt argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var6, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: AVG openb distinct_opt argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var4, Distinct: var3 == DistinctStr, Over: var6}, true
  }  else if id == BIT_AND {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_AND <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_AND openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_AND openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_AND openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == BIT_OR {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_OR <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_OR openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_OR openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_OR openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == BIT_XOR {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_XOR <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_XOR openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_XOR openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: BIT_XOR openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == COUNT {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: COUNT <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.distinctOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: COUNT openb <distinct_opt>', found: 'string(tok)'")
    }
    var4, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: COUNT openb distinct_opt <argument_expression_list>', found: 'string(tok)'")
    }
    var5, tok := p.next()
    if var5 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: COUNT openb distinct_opt argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var6, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: COUNT openb distinct_opt argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var4, Distinct: var3 == DistinctStr, Over: var6}, true
  }  else if id == JSON_ARRAYAGG {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: JSON_ARRAYAGG <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: JSON_ARRAYAGG openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: JSON_ARRAYAGG openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: JSON_ARRAYAGG openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == JSON_OBJECTAGG {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: JSON_OBJECTAGG <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: JSON_OBJECTAGG openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: JSON_OBJECTAGG openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: JSON_OBJECTAGG openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == MIN {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: MIN <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.distinctOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: MIN openb <distinct_opt>', found: 'string(tok)'")
    }
    var4, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: MIN openb distinct_opt <argument_expression_list>', found: 'string(tok)'")
    }
    var5, tok := p.next()
    if var5 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: MIN openb distinct_opt argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var6, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: MIN openb distinct_opt argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var4, Distinct: var3 == DistinctStr, Over: var6}, true
  }  else if id == STDDEV_POP {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV_POP <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV_POP openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV_POP openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV_POP openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == STDDEV {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == STD {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: STD <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: STD openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: STD openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: STD openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == STDDEV_SAMP {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV_SAMP <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV_SAMP openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV_SAMP openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: STDDEV_SAMP openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == SUM {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: SUM <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.distinctOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: SUM openb <distinct_opt>', found: 'string(tok)'")
    }
    var4, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: SUM openb distinct_opt <argument_expression_list>', found: 'string(tok)'")
    }
    var5, tok := p.next()
    if var5 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: SUM openb distinct_opt argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var6, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: SUM openb distinct_opt argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var4, Distinct: var3 == DistinctStr, Over: var6}, true
  }  else if id == VAR_POP {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: VAR_POP <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: VAR_POP openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: VAR_POP openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: VAR_POP openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == VARIANCE {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: VARIANCE <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: VARIANCE openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: VARIANCE openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: VARIANCE openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }  else if id == VAR_SAMP {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_aggregate_with_window: VAR_SAMP <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: VAR_SAMP openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_aggregate_with_window: VAR_SAMP openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    var5, ok := p.overOpt()
    if !ok {
      p.fail("expected: 'function_call_aggregate_with_window: VAR_SAMP openb argument_expression_list closeb <over_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3, Over: var5}, true
  }
  return nil, false
}
func (p *parser) functionCallConflict() (Expr, bool) {
  id, tok := p.peek()
  if id == IF {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_conflict: IF <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_conflict: IF openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_conflict: IF openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == DATABASE {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_conflict: DATABASE <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionListOpt()
    if !ok {
      p.fail("expected: 'function_call_conflict: DATABASE openb <argument_expression_list_opt>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_conflict: DATABASE openb argument_expression_list_opt <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == MOD {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_conflict: MOD <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_conflict: MOD openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_conflict: MOD openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == REPLACE {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_conflict: REPLACE <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_conflict: REPLACE openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_conflict: REPLACE openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == SUBSTR {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_conflict: SUBSTR <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_conflict: SUBSTR openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_conflict: SUBSTR openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: var3}, true
  }  else if id == SUBSTRING {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_conflict: SUBSTRING <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.argumentExpressionList()
    if !ok {
      p.fail("expected: 'function_call_conflict: SUBSTRING openb <argument_expression_list>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != ')' {
      p.fail("expected: 'function_call_conflict: SUBSTRING openb argument_expression_list <')'>', found: '" + string(tok) + "'")
    }
    return &varvar, true
  }
  return nil, false
}
func (p *parser) functionCallNonkeyword() (Expr, bool) {
  id, tok := p.peek()
  if id == CURRENT_DATE {
    var1, _ := p.next()
    var2, ok := p.funcParensOpt()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: CURRENT_DATE <func_parens_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1))}, true
  }  else if id == CURRENT_USER {
    var1, _ := p.next()
    var2, ok := p.funcParensOpt()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: CURRENT_USER <func_parens_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1))}, true
  }  else if id == UTC_DATE {
    var1, _ := p.next()
    var2, ok := p.funcParensOpt()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: UTC_DATE <func_parens_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1))}, true
  }  else if _, ok := p.functionCallOnUpdate(); ok {
    return &var1, true
  }  else if id == CURRENT_TIME {
    var1, _ := p.next()
    var2, ok := p.funcDatetimePrecOpt()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: CURRENT_TIME <func_datetime_prec_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: SelectExprs{&AliasedExpr{Expr: var2}}}, true
  }  else if id == UTC_TIME {
    var1, _ := p.next()
    var2, ok := p.funcDatetimePrecOpt()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: UTC_TIME <func_datetime_prec_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: SelectExprs{&AliasedExpr{Expr: var2}}}, true
  }  else if id == UTC_TIMESTAMP {
    var1, _ := p.next()
    var2, ok := p.funcDatetimePrecOpt()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: UTC_TIMESTAMP <func_datetime_prec_opt>', found: 'string(tok)'")
    }
    return &FuncExpr{Name: NewColIdent(string(var1)), Exprs: SelectExprs{&AliasedExpr{Expr: var2}}}, true
  }  else if id == TIMESTAMPADD {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_nonkeyword: TIMESTAMPADD <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.timeUnit()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: TIMESTAMPADD openb <time_unit>', found: 'string(tok)'")
    }
    var5, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: TIMESTAMPADD openb time_unit ',' <value_expression>', found: 'string(tok)'")
    }
    var7, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: TIMESTAMPADD openb time_unit ',' value_expression ',' <value_expression>', found: 'string(tok)'")
    }
    var8, tok := p.next()
    if var8 != ')' {
      p.fail("expected: 'function_call_nonkeyword: TIMESTAMPADD openb time_unit ',' value_expression ',' value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &TimestampFuncExpr{Name:string("timestampadd"), Unit:string(var3), Expr1:var5, Expr2:var7}, true
  }  else if id == TIMESTAMPDIFF {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_nonkeyword: TIMESTAMPDIFF <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.timeUnit()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: TIMESTAMPDIFF openb <time_unit>', found: 'string(tok)'")
    }
    var5, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: TIMESTAMPDIFF openb time_unit ',' <value_expression>', found: 'string(tok)'")
    }
    var7, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: TIMESTAMPDIFF openb time_unit ',' value_expression ',' <value_expression>', found: 'string(tok)'")
    }
    var8, tok := p.next()
    if var8 != ')' {
      p.fail("expected: 'function_call_nonkeyword: TIMESTAMPDIFF openb time_unit ',' value_expression ',' value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &TimestampFuncExpr{Name:string("timestampdiff"), Unit:string(var3), Expr1:var5, Expr2:var7}, true
  }  else if id == EXTRACT {
    var1, _ := p.next()
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_nonkeyword: EXTRACT <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.timeUnit()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: EXTRACT openb <time_unit>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != FROM {
      p.fail("expected: 'function_call_nonkeyword: EXTRACT openb time_unit <FROM>', found: '" + string(tok) + "'")
    }
    var5, ok := p.valueExpression()
    if !ok {
      p.fail("expected: 'function_call_nonkeyword: EXTRACT openb time_unit FROM <value_expression>', found: 'string(tok)'")
    }
    var6, tok := p.next()
    if var6 != ')' {
      p.fail("expected: 'function_call_nonkeyword: EXTRACT openb time_unit FROM value_expression <')'>', found: '" + string(tok) + "'")
    }
    return &ExtractFuncExpr{Name: string(var1), Unit: string(var3), Expr: var5}, true
  }
  return nil, false
}
func (p *parser) functionCallGeneric() (Expr, bool) {
  id, tok := p.peek()
  if _, ok := p.sqlId(); ok {
    var2, tok := p.next()
    if var2 != '(' {
      p.fail("expected: 'function_call_generic: sql_id <'('>', found: '" + string(tok) + "'")
    }
    var3, ok := p.distinctOpt()
    if !ok {
      p.fail("expected: 'function_call_generic: sql_id openb <distinct_opt>', found: 'string(tok)'")
    }
    var4, ok := p.argumentExpressionListOpt()
    if !ok {
      p.fail("expected: 'function_call_generic: sql_id openb distinct_opt <argument_expression_list_opt>', found: 'string(tok)'")
    }
    var5, tok := p.next()
    if var5 != ')' {
      p.fail("expected: 'function_call_generic: sql_id openb distinct_opt argument_expression_list_opt <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Name: var1, Distinct: var3 == DistinctStr, Exprs: var4}, true
  }  else if _, ok := p.tableId(); ok {
    var3, ok := p.reservedSqlId()
    if !ok {
      p.fail("expected: 'function_call_generic: table_id '.' <reserved_sql_id>', found: 'string(tok)'")
    }
    var4, tok := p.next()
    if var4 != '(' {
      p.fail("expected: 'function_call_generic: table_id '.' reserved_sql_id <'('>', found: '" + string(tok) + "'")
    }
    var5, ok := p.argumentExpressionListOpt()
    if !ok {
      p.fail("expected: 'function_call_generic: table_id '.' reserved_sql_id openb <argument_expression_list_opt>', found: 'string(tok)'")
    }
    var6, tok := p.next()
    if var6 != ')' {
      p.fail("expected: 'function_call_generic: table_id '.' reserved_sql_id openb argument_expression_list_opt <')'>', found: '" + string(tok) + "'")
    }
    return &FuncExpr{Qualifier: var1, Name: var3, Exprs: var5}, true
  }
  return nil, false
}
