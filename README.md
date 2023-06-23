# protothing


```
file: {
  name: "addressbook.proto"
  package: "tutorial"
  dependency: "bq_field.proto"
  message_type: {
    name: "Person"
    field: {
      name: "name"
      number: 1
      label: LABEL_OPTIONAL
      type: TYPE_STRING
      json_name: "name"
    }
    field: {
      name: "email"
      number: 2
      label: LABEL_OPTIONAL
      type: TYPE_STRING
      json_name: "email"
      options: {
        [gen_bq_schema.bigquery]: {
          policy_tags: "pii"
        }
      }
    }
    field: {
      name: "phone"
      number: 3
      label: LABEL_OPTIONAL
      type: TYPE_MESSAGE
      type_name: ".tutorial.Person.PhoneNumber"
      json_name: "phone"
    }
    nested_type: {
      name: "PhoneNumber"
      field: {
        name: "number"
        number: 1
        label: LABEL_OPTIONAL
        type: TYPE_STRING
        json_name: "number"
        options: {
          [gen_bq_schema.bigquery]: {
            policy_tags: "pii"
          }
        }
      }
    }
  }
  message_type: {
    name: "AddressBook"
    field: {
      name: "people"
      number: 1
      label: LABEL_REPEATED
      type: TYPE_MESSAGE
      type_name: ".tutorial.Person"
      json_name: "people"
    }
  }
  syntax: "proto3"
}
```

New PII message that can be added as nested message and field to original message type

field: {
  name: "email"
  number: 2
  label: LABEL_OPTIONAL
  type: TYPE_STRING
  json_name: "email"
}
field: {
  name: "phone"
  number: 3
  label: LABEL_OPTIONAL
  type: TYPE_MESSAGE
  type_name: ".tutorial.Person.PhoneNumberPII"
  json_name: "phone"
}
nested_type: {
  name: "PhoneNumberPII"
  field: {
    name: "number"
    number: 1
    label: LABEL_OPTIONAL
    type: TYPE_STRING
    json_name: "number"
  }
}
```
