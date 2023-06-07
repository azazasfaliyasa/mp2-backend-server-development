# mp2-backend-server-development

# **README**

## Api Implementation - Mini Project 2
1. Subject: Mini Project of Backend Server Development
2. Types:
   ○ Hands-on
3. Topics covered: Go

## Scenario
You are working in an e-commerce company and single handedly responsible with a
CRM ( customer relationship management ) with the user management data. The user
data should synchronize with user data from another system which from
https://reqres.in/api/users?page=2,

## Context
+ For the MVP, the requirements of CRM service are defined below
+ User source data from https://reqres.in/api/users?page=2
+ There is more then one admin and only one super admin
+ The actor who can access the services is admin with role admin and super admin with role super-admin.
+ User role is the customer.

## Milestone 1
+ User management service should be written with Go and using a database from aprevious mini project that you have created.
+ Use gorm, gin-gonic, and library to develop this service.
+ The project should use onion architecture. Follow these template boilerplate onion architecture
+ Initiate the repo using the boilerplates. And there are two module which are:
- Account its consist login, register and CRUD for admin data.
- Customers consist of CRUD to customer data.
+ user management service will expose several APIs with these functions
- An admins and super admin could register a customer at user management services
- An admin could register as admin at user management services
- A Super admin could approve/reject admin registration at user management services
- A Super admin could see approval request at user management services
- An admins and super admin could login at user management services
- An admins and super admin could remove a customer data at user management services
- A super admin could remove a admin data at user management services
- A super admin could activate/deactivate a admin data at user management services
- An admins and super admin could get all a customer data with parameter (search by name and email ) and pagination
- An admins and super admin could get all a admins data with parameter (search by username ) and pagination
- Every time the admin gets a list of customers, service gets data from https://reqres.in/api/users?page=2 and saves into the db if data does not exist.
- A validation is needed to validate the data based on their respective format
- Set repo in gitlab and use the feature branch for development. Merge the feature branch to main when finished work on several features.

## Milestone 2
+ Add a /login api to authenticate users.
+ Password should be hashed
+ Authentication using Basic Auth (username and password)
+ Commit the working code to git

## Milestone 3
+ Add JWT authorization as access_token to access resources after user logged in
+ JWT is passed in Auth header as Bearer token
+ Commit the working code to git
