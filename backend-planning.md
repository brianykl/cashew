VERY IMPORTANT: be careful w ppls financial data, i need to read this and try to enforce best practices
https://ofdss.org/#documents

need-to-dos:
- authorizing users through OAuth2.0
- load monthly-recurring spend and break-down? does this need to be saved in our db or should we grab whenever user logs in? should this be cached?
- monthly recurring spend total, entertainment & media, health & wellness, transportation, food and nutrition, shelter & utilities, etc.
- vendor name, vendor amount, payment type, transaction level data on all of this
- what should users be able to query? chatbot? maybe have gpt or something to ask about how to streamline/prioritize spending  

what to think about:
- how to organize into microservices? need gRPC calls for all microservices
  - user management
    - user registration
    - user authentication
    - profile mgmt
  - transactions analysis
    - plaid api calls to get transactions
    - categorization?
    - pulling/maintaining history
  - account management
    - links new accounts, updating account details, synchronizing through plaid
  - analytics/recommendation engine
    - llm wrapper
  - security/cryptography
    - make stuff safe

api calls to plaid:
- i think most things can be handles with the /transactions endpoint, ntd rn is read through that documentation
- /transactions/sync can subscribe to all transactions associated with an account and get updates like a stream
- /transactions/recurring use after, subscribe to RECURRING_TRANSACTIONS_UPDATE webhook

auth0 flow for plaid:


db schema:
- items we want to track:
  - user
    - user_id (primary key)
    - name <- encrypt at rest
    - email <- encrypt at rest
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
