need-to-dos:
- authorizing users through OAuth2.0
- load monthly-recurring spend and break-down? does this need to be saved in our db or should we grab whenever user logs in? should this be cached?
- monthly recurring spend total, entertainment & media, health & wellness, transportation, food and nutrition, shelter & utilities, etc.
- vendor name, vendor amount, payment type, transaction level data on all of this
- what should users be able to query? chatbot? maybe have gpt or something to ask about how to streamline/prioritize spending  

what to think about:
- how to organize into microservices? need gRPC calls for all microservices
  - authorization
  - api calls to plaid
  - encryption/decryption
  - storing/fetching/deleting data stored in server-side db supporting graphql queries

db schema:
- items we want to track:
  - user
    - user_id (primary key)
    - name
    - email 
    - password <- need to be hashed
  - account
    - account_id (primary key)
    - user_id (foreign key to user)
    - account_type
    - account_name
  - transaction
    - transaction_id (primary key)
    - account_id (foreign key to account)
    - category_id (foreign key to category)
    - amount
    - date
    - vendor_name
    - description
  - category
    - name
    - description
