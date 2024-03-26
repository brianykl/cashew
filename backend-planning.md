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
