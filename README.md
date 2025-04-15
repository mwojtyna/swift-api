# swift-api

Home exercise for potential interns

## Database schema

```mermaid
erDiagram
    bank {
        varchar_11 swift_code PK
        varchar_11 hq_swift_code FK "NOT NULL | INDEX"
        text bank_name "NOT NULL"
        text address  "NOT NULL"
        varchar_2 country_iso2_code "NOT NULL | INDEX"
        text country_name "NOT NULL"
    }
    bank 1--0+ bank: "branches"
```
