# Documentation of API endpoints

API Test Root URL : https://nameless-coast-89815.herokuapp.com/api/v1

Production API Root URL : https://api.tasarruf-admin.com

## Contents

- [User](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#user)

  - [Create User](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#create-user)
  - [Email Login](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#email-login)
  - [Forget Password](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#forget-password)
  - [Is Email Registered](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#is-email-registered)
  - [Is Phone Registered](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#is-phone-registered)
  - [Get User](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-user)
  - [Delete User](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#delete-user)
  - [Update User](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#update-user)
  - [Update Profile Image](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#update-profile-image)
  - [Update Trade License](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#update-trade-license)
  - [Add Partner Photo](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#add-partner-photo)
  - [Delete Partner Photo](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#delete-partner-photo)
  - [Resend Verification Code](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#resend-verification-code)
  - [Verify User](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#verify-user)
  - [Update Main Branch](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#update-main-branch)
  - [Update Password](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#update-password)
  - [Validate Customer Partner Integrity](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#validate-customer-partner-integrity)
  - [Share](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#share)
  - [Search Users](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#search-users)

* [Branch](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#branch)

  - [Create Branch](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#create-branch)
  - [Delete Branch](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#delete-branch)
  - [Edit Branch](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#edit-branch)
  - [Search Branches](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#search-branches)
  - [Get Branch By ID](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-branch-by-id)
  - [Get Branch By Owner](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-branch-by-owner)
  - [Get Branch By Location](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-branch-by-location)
  - [Get My Branches](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-my-branches)

* [Subscription Plans](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#subscription-plans)

  - [Create Subscription Plan](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#create-subscription-plan)
  - [Get All Plans](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-all-plans)
  - [Delete a Plan](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#delete-a-plan)
  - [Update a Plan](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#updated-a-plan)

* [Subscriptions](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#subscriptions)

  - [Subscribe to plan](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#subscribe-to-plan)
  - [Renew Subscription](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#renew-subscription)
  - [Upgrade Subscription](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#upgrade-subscription)
  - [Get My Subscription](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-my-subscription)
  - [Get My Subscription With Partner](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-my-subscription-with-partner)

* [Offers](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#offers)

  - [Consume an offer](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#consume-an-offer)
  - [Connect as a customer](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#connect-as-a-customer)
  - [Get My Offers History](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-my-offers-history)
  - [Send My Offers History Email](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#send-my-offers-history-email)
  - [Get Offer](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-offer)

* [Reviews](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#reviews)

  - [Create Review](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#create-review)
  - [Update Review](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#update-review)
  - [Delete Review](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#delete-review)
  - [Get Review](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-review)
  - [Get My Reviews](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-my-reviews)
  - [Get Partner Reviews](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-partner-reviews)

* [Categories](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#categories)

  - [Create Category](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#create-category)
  - [Edit Category](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#edit-category)
  - [Delete Category](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#delete-category)
  - [Get All Categories](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-all-categories)

* [Cities](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#cities)

  - [Create City](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#create-city)
  - [Delete City](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#delete-city)
  - [Get All Cities](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-all-cities)
  - [Update City](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#update-city)
  - [Get City By ID](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-city-by-id)

* [Support](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#support)

  - [Create Support](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#create-support)
  - [Update Support](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#update-support)
  - [Get Support](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-support)

* [Exclusive Partners](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#exclusive-partners)

  - [Set As Exclusive](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#set-as-exclusive)
  - [Remove As Exclusive](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#remove-as-exclusive)
  - [Get Exclusive Partners](https://github.com/ahmedaabouzied/tasarruf/blob/master/docs/endpoints.md#get-exclusive-partners)

### User

#### Create User

```http
POST /public/user
```

Description : Creates a new user.

The JSON body should have the following parameters.

|      Parameter      |  Type  |             Required             |                                             Validations                                             |
| :-----------------: | :----: | :------------------------------: | :-------------------------------------------------------------------------------------------------: |
|     `firstName`     | string |    true for all account types    |                                     lowercase, 2-50 chars long                                      |
|     `lastName`      | string |    true for all account types    |                                     lowercase, 2-50 chars logn                                      |
|     `brandName`     | string | true for `partner` accounts only |                                                  -                                                  |
|       `phone`       | string | true for `partner` accounts only |                                                  -                                                  |
|      `mobile`       | string |    true for all account types    |                                                  -                                                  |
|       `email`       | string |    true for all account types    |                                                  -                                                  |
|    `accountType`    | string |    true for all account types    |                                either `user` or `pratner` or `admin`                                |
|     `password`      | string |    true for all account types    |                                      must be 8 - 50 chars long                                      |
|      `country`      | string |    true for all account types    |                                 lowercase, must be 2-50 chars long                                  |
|      `cityID`       |  int   |  true for `user` accounts only   |                                                  -                                                  |
| `mainBranchAddress` | string | true for `partner` accounts only |                                                  -                                                  |
|    `dateOfBirth`    | string |  true for `user` accounts only   | Must be in the form of `rfc3339`. More info [here](https://tools.ietf.org/html/rfc3339#section-5.8) |

#### Email Login

```http
POST /public/email-login
```

Description : Login with email.

The JSON body should have the following parameters.

| Parameter  |  Type  | Required | Description |
| :--------: | :----: | :------: | :---------: |
|  `email`   | string |   true   |      -      |
| `password` | string |   true   |      -      |

#### Is Email Registered

```http
GET /public/is-email-registered
```

Description : returns true if the given email is registered on the platform.

Expects `email` in the url query parameters.

#### Is Phone Registered

```http
GET /public/is-phone-registered
```

Description : returns true if the given phone number is registered on the platform.

Expects `phone` in the url query parameters.

#### Forget Password

```http
POST /public/forget-password
```

Description : Sends a password to the user's mobile number. This password expires in 60 minutes.

The JSON body should have the following parameters.

| Parameter |  Type  | Required | Description |
| :-------: | :----: | :------: | :---------: |
| `mobile`  | string |   true   |      -      |

#### Get User

```http
GET /user
```

Description : Returns the currently logged in user

- Headers :
  - Token : {Authentication Token}

#### Delete User

```http
DELETE  /user
```

Description : Deletes the currently logged in user

- Headers :
  - Token : {Authentication Token}

#### Update User

```http
PUT /user
```

Description : Updates user details.

The JSON body should have the following parameters.

|   Parameter   |  Type  |           Required            |                                             Description                                             |
| :-----------: | :----: | :---------------------------: | :-------------------------------------------------------------------------------------------------: |
|  `firstName`  | string |             true              |                                                  -                                                  |
|  `lastName`   | string |             true              |                                                  -                                                  |
|   `country`   | string |             true              |                                                  -                                                  |
|   `cityID`    |  int   |             true              |                                                  -                                                  |
| `dateOfBirth` | string | true for `user` accounts only | Must be in the form of `rfc3339`. More info [here](https://tools.ietf.org/html/rfc3339#section-5.8) |
|   `mobile`    | string |             true              |                                                  -                                                  |
|    `email`    | string |             true              |                                                  -                                                  |

#### Update Profile Image

```http
PUT /user/profile-image
```

Description : Updates the profile image of customer / partner.

- Headers :
  - Token : {Authentication Token}
  - Contet-Type : multipart/form-data

The _multipart_ body should have the following parameters:

|   Parameter    | Type | Required |                  Description                   |
| :------------: | :--: | :------: | :--------------------------------------------: |
| `profileImage` | file |   true   | has .png / .jpeg extension of 3MB maximum size |

#### Update Trade License

```http
PUT /user/trade-license
```

Description : Uploads and updates the partner's trade license.

- Headers :
  - Token : {Authentication Token}
  - Contet-Type : multipart/form-data

The _multipart_ body should have the following parameters:

|   Parameter    | Type | Required |                        Description                        |
| :------------: | :--: | :------: | :-------------------------------------------------------: |
| `tradeLicense` | file |   true   | has .png / .jpeg / .pdf extension with no size constraint |

#### Add Partner Photo

```http
PUT /user/partner-photo
```

Description : Adds a photo to partner's photos.

- Headers :
  - Token : {Authentication Token}
  - Contet-Type : multipart/form-data

The _multipart_ body should have the following parameters:

| Parameter | Type | Required |        Description         |
| :-------: | :--: | :------: | :------------------------: |
|  `photo`  | file |   true   | has .png / .jpeg extension |

#### Delete Partner Photo

```http
DELETE /user/partner-photo/:photoID
```

Description : Deletes the given partner's photo.

- Headers :
  - Token : {Authentication Token}

#### Resend Verification Code

```http
POST /user/resend-verification-code
```

Description : Resends a new verification code to the user's mobile number.

- Headers :
  - Token : {Authentication Token}

#### Verify User

```http
POST /user/verify
```

Description : Sets the user as verified.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters :

| Parameter |  Type  | Required |      Description      |
| :-------: | :----: | :------: | :-------------------: |
|  `code`   | string |   true   | The OTP sent over SMS |

#### Update Main Branch

```http
PUT /user/main-branch
```

Description : Sets the user as verified.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters :

|      Parameter      |  Type  | Required |     Description      |
| :-----------------: | :----: | :------: | :------------------: |
|   `discountValue`   | float  |   true   | offer discount value |
| `mainBranchAddress` | string |   true   |          -           |
|       `phone`       | string |   true   |          -           |
|      `cityID`       |  int   |   true   |          -           |
|      `country`      | string |   true   |          -           |
|     `brandName`     | string |   true   |          -           |
| `offerDiscription`  | string |   true   |          -           |
|    `categroryID`    | string |   true   |          -           |

#### Update Password

```http
POST /user/update-pass
```

Description : Sets the user as verified.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters :

|   Parameter   |  Type  | Required | Description |
| :-----------: | :----: | :------: | :---------: |
| `newPassword` | string |   true   |      -      |

#### Validate Customer Partner Integrity

```http
GET /user/validate-offer?customerID={cusomerID}&partnerID={partnerID}
```

Description : Validates customer partner integriy. Returns the customer, subscription and plan.

- Headers :
  - Token : {Authentication Token}

#### Share

```http
POST /user/share
```

Description : Creates a share record for the currently logged in user. This increases the count of remaining offers of the customer with partners that have the sharing property enabled.

- Headers :
  - Token : {Authentication Token}

#### Search Users

```http
GET /admin/users?q={search term}
```

Description : Searches users by email or phone number.

- Headers :
  - Token : {Authentication Token}

### Branch

#### Create Branch

```http
POST /branch
```

Description : Create a new branch.

The JSON body should have the following parameters.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

| Parameter |  Type  | Required | Description |
| :-------: | :----: | :------: | :---------: |
| `country` | string |   true   |      -      |
| `cityID`  |  int   |   true   |      -      |
| `address` | string |   true   |      -      |
| `mobile`  | string |   true   |      -      |
|  `phone`  | string |   true   |      -      |

#### Edit Branch

```http
PUT /branch/:id
```

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

Description : Edit branch.

The JSON body should have the following parameters.

| Parameter |  Type  | Required | Description |
| :-------: | :----: | :------: | :---------: |
| `country` | string |   true   |      -      |
| `cityID`  |  int   |   true   |      -      |
| `address` | string |   true   |      -      |
| `mobile`  | string |   true   |      -      |
|  `phone`  | string |   true   |      -      |

#### Search Branches

```http
GET /branch?brandName=<brand name>&cityID=<city ID>&categoryID=<category ID>
```

- Headers :
  - Token : {Authentication Token}

Description: Returns branches with the given params.

#### Delete Branch

```http
DELETE /branch/:id
```

- Headers :
  - Token : {Authentication Token}

Description : Deletes branch.

#### Get Branch By ID

```http
GET /branch/:id
```

- Headers :
  - Token : {Authentication Token}

Description : Returns the branch with the given ID.

#### Get Branch By Owner

```http
GET /branches-by-owner/:id
```

- Headers :
  - Token : {Authentication Token}

Description : Returns the branches created by the owner with the given ID.

#### Get Branch By Location

```http
GET /branches-by-location?country={country}&city={cityID}
```

- Headers :
  - Token : {Authentication Token}

Description : Returns the branches within the given location.

Expects country , city as url query parameters.

#### Get My Branches

```http
GET /my-branches
```

- Headers :
  - Token : {Authentication Token}

Description : Returns the branches of the currently logged in user.

### Subscription Plans

#### Create Subscription Plan

```http
POST /plans
```

Description : Create a new subscription plan.

The JSON body should have the following parameters.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

|      Parameter       |  Type   | Required |          Description          |
| :------------------: | :-----: | :------: | :---------------------------: |
|    `englishName`     | string  |   true   | length must be `> 2 and < 50` |
| `englishDescription` | string  |   true   |               -               |
|    `turkishName`     | string  |   true   | length must be `> 2 and < 50` |
|    `turkishName`     | string  |   true   |               -               |
|       `price`        | float64 |   true   |               -               |
|   `countOfOffers`    |  uint   |   true   |               -               |
|       `image`        | string  |   true   |      must be a valid URL      |

#### Get All Plans

```http
GET /plans
```

Description : Returns all the available subscription plans.

- Headers :
  - Token : {Authentication Token}

#### Delete a plan

```http
DELETE /plans/:id
```

Description : Delete the plan with the given ID.

- Headers :
  - Token : {Authentication Token}

#### Update a plan

```http
PUT /plans/:id
```

Description : Updates the plan with the given ID.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

|      Parameter       |  Type   | Required |          Description          |
| :------------------: | :-----: | :------: | :---------------------------: |
|    `englishName`     | string  |   true   | length must be `> 2 and < 50` |
| `englishDescription` | string  |   true   |               -               |
|    `turkishName`     | string  |   true   | length must be `> 2 and < 50` |
|    `turkishName`     | string  |   true   |               -               |
|       `price`        | float64 |   true   |               -               |
|   `countOfOffers`    |  uint   |   true   |               -               |
|       `image`        | string  |   true   |      must be a valid URL      |

### Subscriptions

#### Subscribe to plan

```http
POST subscription/subscribe/:planID
```

Description : Subscribes the currently logged in user to the plan with the given id. This first attempts to create a transaction with the plan price. If transaction is successful then it proceeds with creating the subscription. _The user must be not subscribed to any plan_.

The JSON body should have the following parameters.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

|    Parameter     |  Type  | Required |                                         Description                                         |
| :--------------: | :----: | :------: | :-----------------------------------------------------------------------------------------: |
|   `cardNumber`   | string |   true   |                        15 (AMEX) or 16 (VISA, MC) digits card number                        |
| `cardHolderName` | string |   true   |                     Legal name of the card owner as written on the card                     |
|  `expireMonth`   | string |   true   |                           Expiration month of the card (2 digits)                           |
|   `expireYear`   | string |   true   |                                 Expiration year of the card                                 |
|      `cvc`       | string |   true   |                4 (AMEX) or 3 (VISA, MC, TROY) digits card verification code                 |
|    `idNumber`    | string |   true   | Identity number of buyer. TCKN for Turkish merchants, passport number for foreign merchants |

#### Renew Subscription

```http
POST /subscription/renew
```

Description : Renews the subscription of the currently logged in user with the same plan that is currently subscribed to. This first attempts to create a transaction with the plan price. If transaction is successful then it proceeds with creating the subscription. _The user must be subscribed to a plan_. If the user has any remaining offers left in the previous subscription they are added to the new subscription.

The JSON body should have the following parameters.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

|    Parameter     |  Type  | Required |                                         Description                                         |
| :--------------: | :----: | :------: | :-----------------------------------------------------------------------------------------: |
|   `cardNumber`   | string |   true   |                        15 (AMEX) or 16 (VISA, MC) digits card number                        |
| `cardHolderName` | string |   true   |                     Legal name of the card owner as written on the card                     |
|  `expireMonth`   | string |   true   |                           Expiration month of the card (2 digits)                           |
|   `expireYear`   | string |   true   |                                 Expiration year of the card                                 |
|      `cvc`       | string |   true   |                4 (AMEX) or 3 (VISA, MC, TROY) digits card verification code                 |
|    `idNumber`    | string |   true   | Identity number of buyer. TCKN for Turkish merchants, passport number for foreign merchants |

#### Upgrade Subscription

```http
POST /subscription/upgrade/:planID
```

Description : Renews the subscription of the currently logged in user with the same plan that is currently subscribed to. This first attempts to create a transaction with the plan price. If transaction is successful then it proceeds with creating the subscription. _The user must be subscribed to a plan_. If the user has any remaining offers left in the previous subscription they are added to the new subscription.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters.

|    Parameter     |  Type  | Required |                                         Description                                         |
| :--------------: | :----: | :------: | :-----------------------------------------------------------------------------------------: |
|   `cardNumber`   | string |   true   |                        15 (AMEX) or 16 (VISA, MC) digits card number                        |
| `cardHolderName` | string |   true   |                     Legal name of the card owner as written on the card                     |
|  `expireMonth`   | string |   true   |                           Expiration month of the card (2 digits)                           |
|   `expireYear`   | string |   true   |                                 Expiration year of the card                                 |
|      `cvc`       | string |   true   |                4 (AMEX) or 3 (VISA, MC, TROY) digits card verification code                 |
|    `idNumber`    | string |   true   | Identity number of buyer. TCKN for Turkish merchants, passport number for foreign merchants |

#### Get My Subscription

```http
GET /subscription
```

Description: Returns the subscription of the currently logged in user with the `remainingOffers` value set to the default according to the plan subscribed to.

- Headers :
  - Token : {Authentication Token}

#### Get My Subscription With Partner

```http
GET /subscription/partner/:partnerID
```

Description: Returns the subscription of the currently logged in user with the `remainingOffers` set to the remaining offers count for the given partner.

- Headers :
  - Token : {Authentication Token}

### Offers

#### Consume an offer

```http
POST /offer
```

Description: Used by a _partner_ user. Creates a new offer record between the partner and the customer.The customer remaining count offers gets decremented.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters:

|  Parameter   | Type  | Required |                                                                            Description                                                                            |
| :----------: | :---: | :------: | :---------------------------------------------------------------------------------------------------------------------------------------------------------------: |
|   `amount`   | float |   true   |                                                   The total amount the customer should pay before the discount                                                    |
| `customerID` |  int  |   true   |                                                  ID of the customer (This could be determined by a QR code scan)                                                  |
| `partnerID`  |  int  |   true   | ID of the parnter owning the offer. It's validated against the currently logged in user ID to make sure the customer is offering the QR code to the right partner |

#### Connect as a customer

```http
ws://{rootUrl}/connect
```

Before the partner consumes an offer from a customer, the customer should be connected to the server over [web socket protocol](https://tools.ietf.org/html/rfc6455).If the customer is not connected the partner would receive an error upon consuming the offer and the offer would not be proccessed.

The code below demonstrates connecting to the server over web socket protocol with a javascript client

```js
const rootUrl = <root url here>;
const url = `ws://{rootUrl}/api/v1/connect`;
const c = new WebSocket(url);

// Send the authentication token in the first message after connection.
// The token should be sent in JSON format.
c.onopen = function() {
  c.send(
    JSON.stringify({
      token: "<Authentication Token>"
    })
  );
};
```

#### Get My Offers History

```http
POST /offer/history
```

Description: Retruns the offers related to the currently logged in user (customer / partner). For partners a start and end date must be given as url query paramters to return the offers consumed by this partner within this time range.

- Headers :
  - Token : {Authentication Token}
  - ContentType: application/json

The JSON body should have the following parameters:

|  Parameter  |  Type  | Required |                                             Description                                             |
| :---------: | :----: | :------: | :-------------------------------------------------------------------------------------------------: |
| `startDate` | string |   True   | Must be in the form of `rfc3339`. More info [here](https://tools.ietf.org/html/rfc3339#section-5.8) |
|  `endDate`  | string |   True   | Must be in the form of `rfc3339`. More info [here](https://tools.ietf.org/html/rfc3339#section-5.8) |

#### Send My Offers History Email

```http
POST /offer/history/mail
```

Description: Sends an email to the user with the

- Headers :
  - Token : {Authentication Token}
  - ContentType : application/json

The JSON body should have the following parameters :

|  Parameter  |  Type  | Required |                                             Description                                             |
| :---------: | :----: | :------: | :-------------------------------------------------------------------------------------------------: |
| `startDate` | string |   True   | Must be in the form of `rfc3339`. More info [here](https://tools.ietf.org/html/rfc3339#section-5.8) |
|  `endDate`  | string |   True   | Must be in the form of `rfc3339`. More info [here](https://tools.ietf.org/html/rfc3339#section-5.8) |

#### Get Offer

```http
GET /offer/:offerID
```

Description: Retruns the offer with the given ID.

- Headers :
  - Token : {Authentication Token}

#### Create Review

```http
POST /review
```

Description: Creates a review on a partner by a customer.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters:

|  Parameter   |  Type  | Required |     Description     |
| :----------: | :----: | :------: | :-----------------: |
| `customerID` |  int   |   true   |          -          |
| `partnerID`  |  int   |   true   |          -          |
|   `stars`    |  int   |   true   | must be <= 5 & >= 1 |
|  `content`   | string |  false   | must be <= 5 & >= 1 |

#### Update Review

```http
PUT /review/:reviewID
```

Description: Update a review.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters:

|  Parameter   |  Type  | Required |     Description     |
| :----------: | :----: | :------: | :-----------------: |
| `customerID` |  int   |   true   |          -          |
| `partnerID`  |  int   |   true   |          -          |
|   `stars`    |  int   |   true   | must be <= 5 & >= 1 |
|  `content`   | string |  false   |          -          |

#### Delete Review

```http
DELETE /review/:reviewID
```

Description: Deletes the review with the given ID.

- Headers :
  - Token : {Authentication Token}

#### Get Review

```http
GET /review/:reviewID
```

Description: Gets the review with the given ID.

- Headers :
  - Token : {Authentication Token}

#### Get My Reviews

```http
GET /reviews
```

Description: Gets the reviews of the currently logged in user.

- Headers :
  - Token : {Authentication Token}

#### Get Partner Reviews

```http
GET /reviews/:partnerID
```

Description: Gets the reviews of the partner of the given ID.

- Headers :
  - Token : {Authentication Token}

### Categories

#### Create Category

```http
POST /category
```

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters:

|   Parameter   |  Type  | Required | Description |
| :-----------: | :----: | :------: | :---------: |
| `turkishName` | string |   true   |      -      |
| `englishName` | string |   true   |      -      |

#### Edit Category

```http
PUT /category/:id
```

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters:

|   Parameter   |  Type  | Required | Description |
| :-----------: | :----: | :------: | :---------: |
| `turkishName` | string |   true   |      -      |
| `englishName` | string |   true   |      -      |

#### Delete Category

```http
DELETE /category/:id
```

Description: Deletes the category with the given ID.

- Headers :
  - Token : {Authentication Token}

#### Get All Categories

```http
GET /category
```

Description: Returns all categories.

- Headers :
  - Token : {Authentication Token}

### Cities

#### Create City

```http
POST /city
```

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters:

|   Parameter   |  Type  | Required | Description |
| :-----------: | :----: | :------: | :---------: |
| `turkishName` | string |   true   |      -      |
| `englishName` | string |   true   |      -      |

#### Update City

```http
PUT /city/:id
```

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters:

|   Parameter   |  Type  | Required | Description |
| :-----------: | :----: | :------: | :---------: |
| `turkishName` | string |   true   |      -      |
| `englishName` | string |   true   |      -      |

#### Delete City

```http
DELETE /city/:id
```

Description: Deletes the city with the given ID.

- Headers :
  - Token : {Authentication Token}

#### Get All Cities

```http
GET /public/cities
```

Description: Returns all cities.

- Headers :
  - Token : {Authentication Token}

#### Get City By ID

```http
GET /city/:id
```

Description: Returns all cities.

- Headers :
  - Token : {Authentication Token}

### Support

#### Create Support

```http
POST /support
```

Description : Creates a new support info record.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters:

| Parameter |  Type  | Required | Description |
| :-------: | :----: | :------: | :---------: |
|  `email`  | string |   true   |      -      |
| `mobile`  | string |   true   |      -      |

#### Update Support

```http
PUT /support
```

Description : Updates support info record.

- Headers :
  - Token : {Authentication Token}
  - Content-Type : application/json

The JSON body should have the following parameters:

| Parameter |  Type  | Required | Description |
| :-------: | :----: | :------: | :---------: |
|  `email`  | string |   true   |      -      |
| `mobile`  | string |   true   |      -      |

#### Get Support

```http
GET /support-info
```

Description : Retruns support info record.

- Headers :
  - Token : {Authentication Token}

### Exclusive Partners

#### Set As Exclusive

```http
POST /exclusive/:partnerID
```

Description : Adds partner to exclusive group.

- Headers :
  - Token : {Authentication Token}

#### Remove As Exclusive

```http
DELETE /exclusive/:partnerID
```

Description : Removes partner from exclusive group.

- Headers :
  - Token : {Authentication Token}

#### Get Exclusive Partners

```http
GET /exclusive
```

Description : Returns partners belonging to the exclusive groupd.

- Headers :
  - Token : {Authentication Token}
