%{
package smi

import (
    "strings"
)

%}

%union{
    token Token
    augments string
    description string
    imports []Import
    importIDs []string
    indexes []string
    modules []*parseModule
    object *parseObject
    objects []*parseObject
    objectMap map[string]*parseObject
    orphans []*parseObject
    status Status
    table bool
    val string
    subidentifiers []string
}

%start mibFile

%token ACCESS
%token AGENT_CAPABILITIES
%token APPLICATION
%token AUGMENTS
%token BEGIN
%token BIN_STRING
%token BITS
%token CHOICE
%token COLON_COLON_EQUAL
%token COMMENT
%token CONTACT_INFO
%token CREATION_REQUIRES
%token COUNTER
%token COUNTER32
%token COUNTER64
%token DEFINITIONS
%token DEFVAL
%token DESCRIPTION
%token DISPLAY_HINT
%token DOT_DOT
%token END
%token ENTERPRISE
%token EXPORTS
%token EXTENDS
%token FROM
%token GROUP
%token GAUGE
%token GAUGE32
%token HEX_STRING
%token IDENTIFIER
%token IMPLICIT
%token IMPLIED
%token IMPORTS
%token INCLUDES
%token INDEX
%token INSTALL_ERRORS
%token INTEGER
%token INTEGER32
%token INTEGER64
%token IPADDRESS
%token LAST_UPDATED
%token LOWERCASE_IDENTIFIER
%token MACRO
%token MANDATORY_GROUPS
%token MAX_ACCESS
%token MIN_ACCESS
%token MODULE
%token MODULE_COMPLIANCE
%token MODULE_IDENTITY
%token NEGATIVE_NUMBER
%token NEGATIVE_NUMBER64
%token NOTIFICATION_GROUP
%token NOTIFICATION_TYPE
%token NOTIFICATIONS
%token NUMBER
%token NUMBER64
%token OBJECT
%token OBJECT_GROUP
%token OBJECT_IDENTITY
%token OBJECT_TYPE
%token OBJECTS
%token OCTET
%token OF
%token ORGANIZATION
%token OPAQUE
%token PIB_ACCESS
%token PIB_DEFINITIONS
%token PIB_INDEX
%token PIB_MIN_ACCESS
%token PIB_REFERENCES
%token PIB_TAG
%token POLICY_ACCESS
%token PRODUCT_RELEASE
%token QUOTED_STRING
%token REFERENCE
%token REVISION
%token SEQUENCE
%token SIZE
%token SPECIAL_CHAR
%token STATUS
%token STRING
%token SUBJECT_CATEGORIES
%token SUPPORTS
%token SYNTAX
%token TEXTUAL_CONVENTION
%token TIMETICKS
%token TRAP_TYPE
%token UNIQUENESS
%token UNITS
%token UNIVERSAL
%token UNSIGNED32
%token UNSIGNED64
%token UPPERCASE_IDENTIFIER
%token VALUE
%token VARIABLES
%token VARIATION
%token WRITE_SYNTAX

%%

mibFile : modules
        {
            // Add modules to the module map stored in the lexer
            for _, m := range $$.modules {
                yylex.(*lexer).modules[m.name] = m
            }
            // Clear object, module data from yys
            $$.objectMap = make(map[string]*parseObject)
            $$.orphans = []*parseObject{}
            $$.objects = []*parseObject{}
            $$.modules = []*parseModule{}
        }
        |
        ;

modules : module
        | modules module
        ;

module : moduleName moduleOid definitions COLON_COLON_EQUAL BEGIN exportsClause linkagePart declarationPart END
       {
           m := &parseModule{
               imports: $7.imports,
               name: $1.val,
               objectTree: []*parseObject{},
               orphans: []*parseObject{},
           }
           for _, o := range $8.objects {
               m.objectTree = append(m.objectTree, o)
               o.setModule(m.name)
           }
           for _, o := range $8.orphans {
               m.orphans = append(m.orphans, o)
           }
           $$.addModule(m)
       }
       ;

moduleOid : '{' objectIdentifier '}'
          |
          ;

definitions : DEFINITIONS
            | PIB_DEFINITIONS
            ;

linkagePart : linkageClause
            {
                $$.imports = $1.imports
            }
            |
            {
                $$.imports = nil
            }
            ;

linkageClause : IMPORTS importPart ';'
              {
                  $$.imports = $2.imports
              }
              ;

exportsClause :
              | EXPORTS ';'
              ;

importPart : imports
           {
               $$.imports = $1.imports
           }
           |
           {
               $$.imports = nil
           }
           ;

imports : import
        {
            $$.imports = $1.imports
        }
        | imports import
        {
            $$.imports = append($1.imports, $2.imports...)
        }
        ;

import : importIdentifiers FROM moduleName
       {
           $$.imports = []Import{}
           for _, id := range $1.importIDs {
               $$.imports = append($$.imports,
                   Import{Object: id, Module: $3.token.literal})
           }
       }
       ;

importIdentifiers : importIdentifier
                  {
                      $$.importIDs = []string{$1.token.literal}
                  }
                  | importIdentifiers ',' importIdentifier
                  {
                      $$.importIDs = append($1.importIDs, $3.token.literal)
                  }
                  ;

importIdentifier : LOWERCASE_IDENTIFIER
                 | UPPERCASE_IDENTIFIER
                 | importedKeyword
                 ;

importedKeyword : importedSMIKeyword
                | importedSPPIKeyword
                | BITS
                | INTEGER32
                | IPADDRESS
                | MANDATORY_GROUPS
                | MODULE_COMPLIANCE
                | MODULE_IDENTITY
                | OBJECT_GROUP
                | OBJECT_IDENTITY
                | OBJECT_TYPE
                | OPAQUE
                | TEXTUAL_CONVENTION
                | TIMETICKS
                | UNSIGNED32
                ;

importedSMIKeyword : AGENT_CAPABILITIES
                   | COUNTER32
                   | COUNTER64
                   | GAUGE32
                   | NOTIFICATION_GROUP
                   | NOTIFICATION_TYPE
                   | TRAP_TYPE
                   ;

importedSPPIKeyword : INTEGER64
                    | UNSIGNED64
                    ;

moduleName : UPPERCASE_IDENTIFIER
           {
               $$.val = $1.token.literal
           }
           ;

declarationPart : declarations
                |
                ;

declarations : declaration
             {
                 (&$$).addObject($1.object)
             }
             | declarations declaration
             {
                 (&$$).addObject($2.object)
             }
             ;

declaration : typeDeclaration
            {
                (&$$).setDecl(declTypeAssignment)
            }
            | valueDeclaration
            {
                (&$$).setDecl(declValueAssignment)
            }
            | objectIdentityClause
            {
                (&$$).setDecl(declIdentity)
            }
            | objectTypeClause
            {
                (&$$).setDecl(declObjectType)
            }
            | trapTypeClause
            {
                (&$$).setDecl(declTrapType)
            }
            | notificationTypeClause
            {
                (&$$).setDecl(declNotificationType)
            }
            | moduleIdentityClause
            {
                (&$$).setDecl(declModuleIdentity)
            }
            | moduleComplianceClause
            {
                (&$$).setDecl(declModuleCompliance)
            }
            | objectGroupClause
            {
                (&$$).setDecl(declObjectGroup)
            }
            | notificationGroupClause
            {
                (&$$).setDecl(declNotificationGroup)
            }
            | agentCapabilitiesClause
            {
                (&$$).setDecl(declAgentCapabilities)
            }
            | macroClause
            | error '}'
            ;

macroClause : macroName MACRO END
            ;

macroName : MODULE_IDENTITY
          | OBJECT_TYPE
          | TRAP_TYPE
          | NOTIFICATION_TYPE
          | OBJECT_IDENTITY
          | TEXTUAL_CONVENTION
          | OBJECT_GROUP
          | NOTIFICATION_GROUP
          | MODULE_COMPLIANCE
          | AGENT_CAPABILITIES
          ;

choiceClause : CHOICE '}'
             ;

fuzzyLowercaseIdentifier : LOWERCASE_IDENTIFIER
                         {
                             $$.val = $1.token.literal
                         }
                         | UPPERCASE_IDENTIFIER
                         {
                             $$.val = $1.token.literal
                         }
                         ;

valueDeclaration : fuzzyLowercaseIdentifier OBJECT IDENTIFIER COLON_COLON_EQUAL '{' objectIdentifier '}'
                 {
                    $$.object = &parseObject{
                        object: &Object{
                            Name: $1.val,
                            Oid: strings.Join($6.subidentifiers, "."),
                        },
                    }
                 }
                 ;

typeDeclaration : typeName COLON_COLON_EQUAL typeDeclarationRHS
                ;

typeName : UPPERCASE_IDENTIFIER
         | typeSMI
         | typeSPPIonly
         ;

typeSMI : typeSMIandSPPI
        | typeSMIonly
        ;

typeSMIandSPPI : IPADDRESS
               | TIMETICKS
               | OPAQUE
               | INTEGER32
               | UNSIGNED32
               ;

typeSMIonly : COUNTER32
            | GAUGE32
            | COUNTER64
            ;

typeSPPIonly : INTEGER64
             | UNSIGNED64
             ;

typeDeclarationRHS : Syntax
                   {
                       $$.table = $1.table
                   }
                   | TEXTUAL_CONVENTION DisplayPart STATUS Status DESCRIPTION Text ReferPart SYNTAX Syntax
                   {
                       $$.table = $9.table
                       $$.status = strToStatus($4.val)
                       $$.description = $6.val
                   }
                   | choiceClause
                   ;

conceptualTable : SEQUENCE OF row
                {
                    $$.table = true
                }
                ;

row : UPPERCASE_IDENTIFIER
    ;

entryType : SEQUENCE '{' sequenceItems '}'
          ;

sequenceItems : sequenceItem
              | sequenceItems ',' sequenceItem
              ;

sequenceItem : LOWERCASE_IDENTIFIER sequenceSyntax
             ;

Syntax : ObjectSyntax
       | BITS '{' NamedBits '}'
       ;

sequenceSyntax : sequenceObjectSyntax
               | BITS
               | UPPERCASE_IDENTIFIER anySubType
               ;

NamedBits : NamedBit
          | NamedBits ',' NamedBit
          ;

NamedBit : LOWERCASE_IDENTIFIER '(' NUMBER ')'
         ;

objectIdentityClause : LOWERCASE_IDENTIFIER OBJECT_IDENTITY STATUS Status DESCRIPTION Text ReferPart COLON_COLON_EQUAL '{' objectIdentifier '}'
                     {
                         $$.object = &parseObject{
                             object: &Object{
                                 Description: $6.val,
                                 Name: $1.token.literal,
                                 Oid: strings.Join($10.subidentifiers, "."),
                                 Status: strToStatus($4.val),
                             },
                             decl: declIdentity,
                         }
                     }
                     ;

objectTypeClause : LOWERCASE_IDENTIFIER OBJECT_TYPE SYNTAX Syntax UnitsPart MaxOrPIBAccessPart SPPIPibReferencesPart SPPIPibTagPart STATUS Status descriptionClause SPPIErrorsPart ReferPart IndexPart MibIndex SPPIUniquePart DefValPart COLON_COLON_EQUAL '{' ObjectName '}'
                 {
                     $$.object = &parseObject{
                         object: &Object{
                             Access: strToAccess($6.val),
                             Description: $11.val,
                             Indexes: $15.indexes,
                             Name: $1.token.literal,
                             Oid: strings.Join($20.subidentifiers, "."),
                             Status: strToStatus($10.val),
                         },
                         decl: declObjectType,
                         table: $4.table,
                         augments: $14.augments,
                     }
                 }
                 ;

descriptionClause : DESCRIPTION Text
                  {
                      $$.val = $2.val
                  }
                  |
                  ;

trapTypeClause : fuzzyLowercaseIdentifier TRAP_TYPE ENTERPRISE objectIdentifier VarPart DescrPart ReferPart COLON_COLON_EQUAL NUMBER
               ;

VarPart : VARIABLES '{' VarTypes '}'
        |
        ;

VarTypes : VarType
         | VarTypes ',' VarType
         ;

VarType : ObjectName
        ;

DescrPart : DESCRIPTION Text
          |
          ;

MaxOrPIBAccessPart : MaxAccessPart
                   | PibAccessPart
                   |
                   ;

PibAccessPart : PibAccess Access
              ;

PibAccess : POLICY_ACCESS
          | PIB_ACCESS
          ;

SPPIPibReferencesPart : PIB_REFERENCES
                      |
                      ;

SPPIPibTagPart : PIB_TAG
               |
               ;

SPPIUniquePart : UNIQUENESS '{' UniqueTypesPart '}'
               |
               ;

UniqueTypesPart : UniqueTypes
                |
                ;

UniqueTypes : UniqueType
            | UniqueTypes ',' UniqueType
            ;

UniqueType : ObjectName
           ;

SPPIErrorsPart : INSTALL_ERRORS '{' Errors '}'
               |
               ;

Errors : Error
       | Errors ',' Error
       ;

Error : LOWERCASE_IDENTIFIER '(' NUMBER ')'
      ;

MaxAccessPart : MAX_ACCESS Access
              {
                  $$.val = $2.token.literal
              }
              | ACCESS Access
              {
                  $$.val = $2.token.literal
              }
              ;

notificationTypeClause : LOWERCASE_IDENTIFIER NOTIFICATION_TYPE NotificationObjectsPart STATUS Status DESCRIPTION Text ReferPart COLON_COLON_EQUAL '{' NotificationName '}'
                       ;


moduleIdentityClause : LOWERCASE_IDENTIFIER MODULE_IDENTITY SubjectCategoriesPart LAST_UPDATED ExtUTCTime ORGANIZATION Text CONTACT_INFO Text DESCRIPTION Text RevisionPart COLON_COLON_EQUAL '{' objectIdentifier '}'
                     {
                         $$.object = &parseObject{
                             object: &Object{
                                 Name: $1.token.literal,
                                 Oid: strings.Join($15.subidentifiers, "."),
                                 Description: $11.val,
                             },
                         }
                     }
                     ;

SubjectCategoriesPart : SUBJECT_CATEGORIES '{' SubjectCategories '}'
                      |
                      ;

SubjectCategories : CategoryIDs
                  ;

CategoryIDs : CategoryID
            | CategoryIDs ',' CategoryID
            ;

CategoryID : LOWERCASE_IDENTIFIER
           | LOWERCASE_IDENTIFIER '(' NUMBER ')'
           ;

ObjectSyntax : SimpleSyntax
             | typeTag SimpleSyntax
             | conceptualTable
             | row
             | entryType
             | ApplicationSyntax
             ;

typeTag : '[' APPLICATION NUMBER ']' IMPLICIT
        | '[' UNIVERSAL NUMBER ']' IMPLICIT
        ;

sequenceObjectSyntax : sequenceSimpleSyntax
                     | sequenceApplicationSyntax
                     ;

valueofObjectSyntax : valueofSimpleSyntax
                    ;

SimpleSyntax : INTEGER
             | INTEGER integerSubType
             | INTEGER enumSpec
             | INTEGER32
             | INTEGER32 integerSubType
             | UPPERCASE_IDENTIFIER enumSpec
             | moduleName '.' UPPERCASE_IDENTIFIER enumSpec
             | UPPERCASE_IDENTIFIER integerSubType
             | moduleName '.' UPPERCASE_IDENTIFIER integerSubType
             | OCTET STRING
             | OCTET STRING octetStringSubType
             | UPPERCASE_IDENTIFIER octetStringSubType
             | moduleName '.' UPPERCASE_IDENTIFIER octetStringSubType
	         | OBJECT IDENTIFIER anySubType
             ;

valueofSimpleSyntax : NUMBER
                    | NEGATIVE_NUMBER
                    | NUMBER64
                    | NEGATIVE_NUMBER64
                    | BIN_STRING
                    | HEX_STRING
                    | LOWERCASE_IDENTIFIER
                    | QUOTED_STRING
                    | '{' objectIdentifier_defval '}'
                    ;

sequenceSimpleSyntax : INTEGER anySubType
                     | INTEGER32 anySubType
                     | OCTET STRING anySubType
                     | OBJECT IDENTIFIER anySubType
                     ;

ApplicationSyntax : IPADDRESS anySubType
                  | COUNTER32 anySubType
                  | GAUGE32
                  | GAUGE32 integerSubType
                  | UNSIGNED32
                  | UNSIGNED32 integerSubType
                  | TIMETICKS anySubType
                  | OPAQUE
                  | OPAQUE octetStringSubType
                  | COUNTER64 anySubType
                  | INTEGER64
                  | INTEGER64 integerSubType
                  | UNSIGNED64
                  | UNSIGNED64 integerSubType
                  ;

sequenceApplicationSyntax : IPADDRESS anySubType
                          | COUNTER32 anySubType
                          | GAUGE32 anySubType
                          | UNSIGNED32 anySubType
                          | TIMETICKS anySubType
                          | OPAQUE
                          | COUNTER64 anySubType
                          | INTEGER64
                          | UNSIGNED64
                          ;

anySubType : integerSubType
           | octetStringSubType
           | enumSpec
           |
           ;

integerSubType : '(' ranges ')'
               ;

octetStringSubType : '(' SIZE '(' ranges ')' ')'
                   ;

ranges : range
       | ranges '|' range
       ;

range : value
      | value DOT_DOT value
      ;

value : NEGATIVE_NUMBER
      | NUMBER
      | NEGATIVE_NUMBER64
      | NUMBER64
      | HEX_STRING
      | BIN_STRING
      ;

enumSpec : '{' enumItems '}'
         ;

enumItems : enumItem
          | enumItems ',' enumItem
          ;

enumItem : LOWERCASE_IDENTIFIER '(' enumNumber ')'
         ;

enumNumber : NUMBER
           | NEGATIVE_NUMBER
           ;

Status : LOWERCASE_IDENTIFIER
       {
           $$.val = $1.token.literal
       }
       ;

Status_Capabilities : LOWERCASE_IDENTIFIER
                    ;

DisplayPart : DISPLAY_HINT Text
            |
            ;

UnitsPart : UNITS Text
          |
          ;

Access : LOWERCASE_IDENTIFIER
       ;

IndexPart : PIB_INDEX '{' Entry '}'
          {
             $$.augments = ""
          }
          | AUGMENTS '{' Entry '}'
          {
             $$.augments = $3.subidentifiers[0]
          }
          | EXTENDS '{' Entry '}'
          {
             $$.augments = ""
          }
          |
          {
             $$.augments = ""
          }
          ;

MibIndex : INDEX '{' IndexTypes '}'
         {
             $$.indexes = $3.indexes
         }
         |
         {
             $$.indexes = nil
         }
         ;

IndexTypes : IndexType
           {
               if $1.val != "" {
                   $$.indexes = []string{$1.val}
               }
           }
           | IndexTypes ',' IndexType
           {
               if $3.val != "" {
                   $$.indexes = append($1.indexes, $3.val)
               }
           }
           ;

IndexType : IMPLIED Index
          | Index
          {
              $$.val = strings.Join($1.subidentifiers, " ")
          }
          ;

Index : ObjectName
      ;

Entry : ObjectName
      ;

DefValPart : DEFVAL '{' Value '}'
           |
           ;

Value : valueofObjectSyntax
      | '{' BitsValue '}'
      ;

BitsValue : BitNames
          |
          ;

BitNames : LOWERCASE_IDENTIFIER
         | BitNames ',' LOWERCASE_IDENTIFIER
         ;

ObjectName : objectIdentifier
            ;

NotificationName : objectIdentifier
                 ;

ReferPart : REFERENCE Text
          |
          ;

RevisionPart : Revisions
             |
             ;

Revisions : Revision
          | Revisions Revision
          ;

Revision : REVISION ExtUTCTime DESCRIPTION Text
         ;

NotificationObjectsPart : OBJECTS '{' Objects '}'
                        |
                        ;

ObjectGroupObjectsPart : OBJECTS '{' Objects '}'
                       ;

Objects : Object
        | Objects ',' Object
        ;

Object : ObjectName
       ;

NotificationsPart : NOTIFICATIONS '{' Notifications '}'
                  ;

Notifications : Notification
              | Notifications ',' Notification
              ;

Notification : NotificationName
             ;

Text : QUOTED_STRING
     {
         $$.val = $1.token.literal
     }
     ;

ExtUTCTime : QUOTED_STRING
           ;

objectIdentifier : subidentifiers
                 ;

subidentifiers : subidentifier
               {
                   $$.subidentifiers = []string{$1.val}
               }
               | subidentifiers subidentifier
               {
                   $$.subidentifiers = append($1.subidentifiers, $2.val)
               }
               ;

subidentifier : fuzzyLowercaseIdentifier
              {
                  $$.val = $1.token.literal
              }
              | moduleName '.' LOWERCASE_IDENTIFIER
              {
                  $$.val = $1.token.literal
              }
              | NUMBER
              {
                  $$.val = $1.token.literal
              }
              | LOWERCASE_IDENTIFIER '(' NUMBER ')'
              {
                  $$.val = $3.token.literal
              }
              | moduleName '.' LOWERCASE_IDENTIFIER '(' NUMBER ')'
              {
                  $$.val = $1.token.literal
              }
              ;

objectIdentifier_defval : subidentifiers_defval
                        ;

subidentifiers_defval : subidentifier_defval
                      | subidentifiers_defval subidentifier_defval
                      ;

subidentifier_defval : LOWERCASE_IDENTIFIER '(' NUMBER ')'
                     | NUMBER
                     ;

objectGroupClause : LOWERCASE_IDENTIFIER OBJECT_GROUP ObjectGroupObjectsPart STATUS Status DESCRIPTION Text ReferPart COLON_COLON_EQUAL '{' objectIdentifier '}'
                  {
                      // XXX TODO
                  }
                  ;

notificationGroupClause : LOWERCASE_IDENTIFIER NOTIFICATION_GROUP NotificationsPart STATUS Status DESCRIPTION Text ReferPart COLON_COLON_EQUAL '{' objectIdentifier '}'
                        {
                            // XXX TODO
                        }
                        ;

moduleComplianceClause : LOWERCASE_IDENTIFIER MODULE_COMPLIANCE STATUS Status DESCRIPTION Text ReferPart ComplianceModulePart COLON_COLON_EQUAL '{' objectIdentifier '}'
                       {
                           /// XXX TODO
                       }
                       ;

ComplianceModulePart : ComplianceModules
                     ;

ComplianceModules : ComplianceModule
                  | ComplianceModules ComplianceModule
                  ;

ComplianceModule : MODULE ComplianceModuleName MandatoryPart CompliancePart
                 ;

ComplianceModuleName : UPPERCASE_IDENTIFIER objectIdentifier
                     | UPPERCASE_IDENTIFIER
                     |
                     ;

MandatoryPart : MANDATORY_GROUPS '{' MandatoryGroups '}'
              |
              ;

MandatoryGroups : MandatoryGroup
                | MandatoryGroups ',' MandatoryGroup
                ;

MandatoryGroup : objectIdentifier
               ;

CompliancePart : Compliances
               |
               ;

Compliances : Compliance
            | Compliances Compliance
            ;

Compliance : ComplianceGroup
           | ComplianceObject
           ;

ComplianceGroup : GROUP objectIdentifier DESCRIPTION Text
                ;

ComplianceObject : OBJECT ObjectName SyntaxPart WriteSyntaxPart AccessPart DESCRIPTION Text
                 ;

SyntaxPart : SYNTAX Syntax
           |
           ;

WriteSyntaxPart : WRITE_SYNTAX WriteSyntax
                |
                ;

WriteSyntax : Syntax
            ;

AccessPart : MIN_ACCESS Access
           | PIB_MIN_ACCESS Access
           |
           ;

agentCapabilitiesClause : LOWERCASE_IDENTIFIER AGENT_CAPABILITIES PRODUCT_RELEASE Text STATUS Status_Capabilities DESCRIPTION Text ReferPart ModulePart_Capabilities COLON_COLON_EQUAL '{' objectIdentifier '}'
                        {
                            // XXX TODO
                        }
                        ;

ModulePart_Capabilities : Modules_Capabilities
                        |
                        ;

Modules_Capabilities : Module_Capabilities
                     | Modules_Capabilities Module_Capabilities
                     ;

Module_Capabilities : SUPPORTS ModuleName_Capabilities INCLUDES '{' CapabilitiesGroups '}' VariationPart
                    ;

CapabilitiesGroups : CapabilitiesGroup
                   | CapabilitiesGroups ',' CapabilitiesGroup
                   ;

CapabilitiesGroup : objectIdentifier
                  ;

ModuleName_Capabilities : UPPERCASE_IDENTIFIER objectIdentifier
                        | UPPERCASE_IDENTIFIER
                        ;

VariationPart : Variations
              |
              ;

Variations : Variation
           | Variations Variation
           ;

Variation : VARIATION ObjectName SyntaxPart WriteSyntaxPart VariationAccessPart CreationPart DefValPart DESCRIPTION Text
          ;

VariationAccessPart : ACCESS VariationAccess
                    |
                    ;

VariationAccess : LOWERCASE_IDENTIFIER
                ;

CreationPart : CREATION_REQUIRES '{' Cells '}'
             |
             ;

Cells : Cell
      | Cells ',' Cell
      ;

Cell : ObjectName
     ;

%%
