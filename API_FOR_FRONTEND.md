Backend API Contract (Front-end)

Overview
- Base URL: server root; paths below are relative to '/'.
- All JSON requests must be valid JSON and must not contain unknown fields (unknown fields are rejected).
- Dates/times are RFC3339 strings (Go time.Time JSON format), e.g. "2024-01-02T15:04:05Z".
- Error responses use the shared format shown in "Error Response".
- Note: several routes are spelled with the prefix "emplpoyee" (typo is in the server routes and must be used as-is).

Auth & Headers
- Authorization: Bearer <accessToken> for protected endpoints (see below).
- Content-Type: application/json for JSON bodies.
- Accept-Language: optional; expected values "tg", "ru", "en". If not set, backend defaults to "ru" for language-sensitive operations.
  - LanguageCode in requests is usually ignored/overwritten by server for most "Create" endpoints and for details update.

Error Response (JSON)
{
  "error": "message",
  "code": "VALIDATION|NOT_FOUND|CONFLICT|UNAUTHORIZED|FORBIDDEN|UNAVAILABLE|INTERNAL",
  "fields": {
    "fieldName": "error message"
  }
}

Health
- GET /health
  - Auth: no
  - Response: 200 text/plain "ok"

Auth
- POST /auth/register
  - Auth: no
  - Body: RegisterRequest
  - Response: 201 (empty body)

- POST /auth/login
  - Auth: no
  - Body: AuthRequest
  - Response: 200 AuthResponse

- POST /auth/refresh
  - Auth: no
  - Body: RefreshTokenRequest
  - Response: 200 AuthResponse

- POST /auth/logout
  - Auth: no
  - Body: LogoutRequest
  - Response: 204 (empty body)

Employee
- GET /employee
  - Auth: no
  - Query: uniqueID (string, required)
  - Response: 200 EmployeeResponse

- GET /employee/profile-picture/{uid}
  - Auth: no
  - Path: uid (string, required)
  - Response: 200 file content (image or binary)

- PUT /employee/profile-picture/{uid}
  - Auth: yes (Bearer)
  - Path: uid (string, required)
  - Body: multipart/form-data with field "file"; max size 10MB
  - Response: 204 (empty body)

Employee Details
- GET /emplpoyee-detail
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeeDetailsResponse[]

- PUT /emplpoyee-detail
  - Auth: yes (Bearer)
  - Body: UpdateFullEmployeeData
  - Response: 200 EmployeeDetailsResponse[] (created/updated items)
  - Note: languageCode for each item is overwritten from Accept-Language.

Employee Degree
- GET /emplpoyee-degree
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeeDegreeResponse[]

- POST /emplpoyee-degree
  - Auth: yes (Bearer)
  - Body: CreateEmployeeDegreeRequest
  - Response: 201 EmployeeDegreeResponse
  - Note: languageCode is derived from Accept-Language.

- PUT /emplpoyee-degree
  - Auth: yes (Bearer)
  - Body: UpdateEmployeeDegreeRequest
  - Response: 200 EmployeeDegreeResponse

- DELETE /emplpoyee-degree/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Employee Main Research Area
- GET /emplpoyee-main-research-area
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeeMainResearchAreaResponse[]

- POST /emplpoyee-main-research-area
  - Auth: yes (Bearer)
  - Body: CreateEmployeeMainResearchAreaRequest
  - Response: 201 EmployeeMainResearchAreaResponse
  - Note: languageCode for main request + keyTopics is overwritten from Accept-Language.

- PUT /emplpoyee-main-research-area
  - Auth: yes (Bearer)
  - Body: UpdateEmployeeMainResearchAreaRequest
  - Response: 200 EmployeeMainResearchAreaResponse

- DELETE /emplpoyee-main-research-area/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Employee Participation In Event
- GET /emplpoyee-participation-in-event
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeeParticipationInEventResponse[]

- POST /emplpoyee-participation-in-event
  - Auth: yes (Bearer)
  - Body: CreateEmployeeParticipationInEventRequest
  - Response: 201 EmployeeParticipationInEventResponse
  - Note: languageCode is derived from Accept-Language.

- PUT /emplpoyee-participation-in-event
  - Auth: yes (Bearer)
  - Body: UpdateEmployeeParticipationInEventRequest
  - Response: 200 EmployeeParticipationInEventResponse

- DELETE /emplpoyee-participation-in-event/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Employee Participation In Professional Community
- GET /emplpoyee-participation-in-professional-community
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeeParticipationInProfessionalCommunityResponse[]

- POST /emplpoyee-participation-in-professional-community
  - Auth: yes (Bearer)
  - Body: CreateEmployeeParticipationInProfessionalCommunityRequest
  - Response: 201 EmployeeParticipationInProfessionalCommunityResponse
  - Note: languageCode is derived from Accept-Language.

- PUT /emplpoyee-participation-in-professional-community
  - Auth: yes (Bearer)
  - Body: UpdateEmployeeParticipationInProfessionalCommunityRequest
  - Response: 200 EmployeeParticipationInProfessionalCommunityResponse

- DELETE /emplpoyee-participation-in-professional-community/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Employee Patent
- GET /emplpoyee-patent
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeePatentResponse[]

- POST /emplpoyee-patent
  - Auth: yes (Bearer)
  - Body: CreateEmployeePatentRequest
  - Response: 201 EmployeePatentResponse
  - Note: languageCode is derived from Accept-Language.

- PUT /emplpoyee-patent
  - Auth: yes (Bearer)
  - Body: UpdateEmployeePatentRequest
  - Response: 200 EmployeePatentResponse

- DELETE /emplpoyee-patent/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Employee Publication
- GET /emplpoyee-publication
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeePublicationResponse[]

- POST /emplpoyee-publication
  - Auth: yes (Bearer)
  - Body: CreateEmployeePublicationRequest
  - Response: 201 EmployeePublicationResponse
  - Note: languageCode is derived from Accept-Language.

- PUT /emplpoyee-publication
  - Auth: yes (Bearer)
  - Body: UpdateEmployeePublicationRequest
  - Response: 200 EmployeePublicationResponse

- DELETE /emplpoyee-publication/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Employee Refresher Course
- GET /emplpoyee-refresher-course
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeeRefresherCourseResponse[]

- POST /emplpoyee-refresher-course
  - Auth: yes (Bearer)
  - Body: CreateEmployeeRefresherCourseRequest
  - Response: 201 EmployeeRefresherCourseResponse
  - Note: languageCode is derived from Accept-Language.

- PUT /emplpoyee-refresher-course
  - Auth: yes (Bearer)
  - Body: UpdateEmployeeRefresherCourseRequest
  - Response: 200 EmployeeRefresherCourseResponse

- DELETE /emplpoyee-refresher-course/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Employee Research Activity
- GET /emplpoyee-research-activity
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeeResearchActivityResponse[]

- POST /emplpoyee-research-activity
  - Auth: yes (Bearer)
  - Body: CreateEmployeeResearchActivityRequest
  - Response: 201 EmployeeResearchActivityResponse
  - Note: languageCode is derived from Accept-Language.

- PUT /emplpoyee-research-activity
  - Auth: yes (Bearer)
  - Body: UpdateEmployeeResearchActivityRequest
  - Response: 200 EmployeeResearchActivityResponse

- DELETE /emplpoyee-research-activity/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Employee Scientific Award
- GET /emplpoyee-scientific-award
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeeScientificAwardResponse[]

- POST /emplpoyee-scientific-award
  - Auth: yes (Bearer)
  - Body: CreateEmployeeScientificAwardRequest
  - Response: 201 EmployeeScientificAwardResponse
  - Note: languageCode is derived from Accept-Language.

- PUT /emplpoyee-scientific-award
  - Auth: yes (Bearer)
  - Body: UpdateEmployeeScientificAwardRequest
  - Response: 200 EmployeeScientificAwardResponse

- DELETE /emplpoyee-scientific-award/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Employee Social
- GET /emplpoyee-social
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeeSocialResponse[]

- POST /emplpoyee-social
  - Auth: yes (Bearer)
  - Body: CreateEmployeeSocialRequest
  - Response: 201 EmployeeSocialResponse

- PUT /emplpoyee-social
  - Auth: yes (Bearer)
  - Body: UpdateEmployeeSocialRequest
  - Response: 200 EmployeeSocialResponse

- DELETE /emplpoyee-social/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Employee Work Experience
- GET /emplpoyee-work-experience
  - Auth: no
  - Query: employeeID (int64, required)
  - Response: 200 EmployeeWorkExperienceResponse[]

- POST /emplpoyee-work-experience
  - Auth: yes (Bearer)
  - Body: CreateEmployeeWorkExperienceRequest
  - Response: 201 EmployeeWorkExperienceResponse
  - Note: languageCode is derived from Accept-Language.

- PUT /emplpoyee-work-experience
  - Auth: yes (Bearer)
  - Body: UpdateEmployeeWorkExperienceRequest
  - Response: 200 EmployeeWorkExperienceResponse

- DELETE /emplpoyee-work-experience/{id}
  - Auth: yes (Bearer)
  - Path: id (int64, required)
  - Response: 204 (empty body)

Request/Response DTOs
RegisterRequest
{
  "tin": "string",
  "gender": "string",
  "email": "string",
  "password": "string"
}

AuthRequest
{
  "email": "string",
  "password": "string"
}

AuthResponse
{
  "accessToken": "string",
  "refreshToken": "string",
  "userID": "string"
}

LogoutRequest
{
  "refreshToken": "string"
}

RefreshTokenRequest
{
  "refreshToken": "string"
}

EmployeeResponse
{
  "id": 0,
  "uniqueID": "string",
  "gender": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339",
  "details": [EmployeeDetailsResponse],
  "degrees": [EmployeeDegreeResponse],
  "workExperiences": [EmployeeWorkExperienceResponse],
  "mainResearchAreas": [EmployeeMainResearchAreaResponse],
  "publications": [EmployeePublicationResponse],
  "scientificAwards": [EmployeeScientificAwardResponse],
  "patents": [EmployeePatentResponse],
  "participationInProfessionalCommunities": [EmployeeParticipationInProfessionalCommunityResponse],
  "refresherCourses": [EmployeeRefresherCourseResponse],
  "participationInEvents": [EmployeeParticipationInEventResponse],
  "researchActivities": [EmployeeResearchActivityResponse],
  "socials": [EmployeeSocialResponse]
}

EmployeeDetailsResponse
{
  "id": 0,
  "languageCode": "string",
  "surname": "string",
  "name": "string",
  "middlename": "string",
  "isNewEmployeeDetails": true,
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

UpdateFullEmployeeData
{
  "data": [UpdateEmployeeDetailsRequest]
}

UpdateEmployeeDetailsRequest
{
  "id": 0,
  "employeeID": 0,
  "languageCode": "string",
  "surname": "string|null",
  "name": "string|null",
  "middlename": "string|null",
  "isNewEmployeeDetails": true
}

CreateEmployeeDegreeRequest
{
  "employeeID": 0,
  "rfInstitutionID": 0,
  "degreeLevel": "string",
  "institutionName": "string",
  "speciality": "string",
  "dateStart": "RFC3339",
  "dateEnd": "RFC3339",
  "givenBy": "string",
  "dateDegreeRecieved": "RFC3339"
}

UpdateEmployeeDegreeRequest
{
  "id": 0,
  "rfInstitutionID": "number|null",
  "degreeLevel": "string|null",
  "institutionName": "string|null",
  "speciality": "string|null",
  "dateStart": "RFC3339|null",
  "dateEnd": "RFC3339|null",
  "givenBy": "string|null",
  "dateDegreeRecieved": "RFC3339|null"
}

EmployeeDegreeResponse
{
  "id": 0,
  "rfInstitutionID": 0,
  "degreeLevel": "string",
  "institutionName": "string",
  "speciality": "string",
  "dateStart": "RFC3339",
  "dateEnd": "RFC3339",
  "givenBy": "string",
  "dateDegreeRecieved": "RFC3339",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

CreateEmployeeMainResearchAreaRequest
{
  "employeeID": 0,
  "languageCode": "string",
  "area": "string",
  "discipline": "string",
  "keyTopics": [CreateResearchAreaKeyTopicRequest]
}

CreateResearchAreaKeyTopicRequest
{
  "languageCode": "string",
  "keyTopicTitle": "string"
}

UpdateEmployeeMainResearchAreaRequest
{
  "id": 0,
  "discipline": "string|null",
  "area": "string|null",
  "keyTopics": [UpdateResearchAreaKeyTopicRequest]
}

UpdateResearchAreaKeyTopicRequest
{
  "id": 0,
  "keyTopicTitle": "string|null"
}

EmployeeMainResearchAreaResponse
{
  "id": 0,
  "discipline": "string",
  "area": "string",
  "keyTopics": [ResearchAreaKeyTopicResponse],
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

ResearchAreaKeyTopicResponse
{
  "id": 0,
  "keyTopicTitle": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

CreateEmployeeParticipationInEventRequest
{
  "employeeID": 0,
  "eventTitle": "string",
  "eventDate": "RFC3339"
}

UpdateEmployeeParticipationInEventRequest
{
  "id": 0,
  "eventTitle": "string|null",
  "eventDate": "RFC3339|null"
}

EmployeeParticipationInEventResponse
{
  "id": 0,
  "eventTitle": "string",
  "eventDate": "RFC3339",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

CreateEmployeeParticipationInProfessionalCommunityRequest
{
  "employeeID": 0,
  "professionalCommunityTitle": "string",
  "roleInProfessionalCommunity": "string"
}

UpdateEmployeeParticipationInProfessionalCommunityRequest
{
  "id": 0,
  "professionalCommunityTitle": "string|null",
  "roleInProfessionalCommunity": "string|null"
}

EmployeeParticipationInProfessionalCommunityResponse
{
  "id": 0,
  "professionalCommunityTitle": "string",
  "roleInProfessionalCommunity": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

CreateEmployeePatentRequest
{
  "employeeID": 0,
  "patentTitle": "string",
  "description": "string"
}

UpdateEmployeePatentRequest
{
  "id": 0,
  "patentTitle": "string|null",
  "description": "string|null"
}

EmployeePatentResponse
{
  "id": 0,
  "patentTitle": "string",
  "description": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

CreateEmployeePublicationRequest
{
  "employeeID": 0,
  "rfPublicationTypeID": "number",
  "name": "string",
  "type": "string",
  "authors": "string|null",
  "journalName": "string|null",
  "volume": "string|null",
  "number": "string|null",
  "pages": "string|null",
  "year": "number|null",
  "link": "string"
}

UpdateEmployeePublicationRequest
{
  "id": 0,
  "rfPublicationTypeID": "number|null",
  "name": "string|null",
  "type": "string|null",
  "authors": "string|null",
  "journalName": "string|null",
  "volume": "string|null",
  "number": "string|null",
  "pages": "string|null",
  "year": "number|null",
  "link": "string|null"
}

EmployeePublicationResponse
{
  "id": 0,
  "rfPublicationTypeID": "number",
  "name": "string",
  "type": "string",
  "authors": "string|null",
  "journalName": "string|null",
  "volume": "string|null",
  "number": "string|null",
  "pages": "string|null",
  "year": "number|null",
  "link": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

CreateEmployeeRefresherCourseRequest
{
  "employeeID": 0,
  "courseTitle": "string",
  "dateStart": "RFC3339",
  "dateEnd": "RFC3339"
}

UpdateEmployeeRefresherCourseRequest
{
  "id": 0,
  "courseTitle": "string|null",
  "dateStart": "RFC3339|null",
  "dateEnd": "RFC3339|null"
}

EmployeeRefresherCourseResponse
{
  "id": 0,
  "courseTitle": "string",
  "dateStart": "RFC3339",
  "dateEnd": "RFC3339",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

CreateEmployeeResearchActivityRequest
{
  "employeeID": 0,
  "researchActivityTitle": "string",
  "employeeRole": "string"
}

UpdateEmployeeResearchActivityRequest
{
  "id": 0,
  "researchActivityTitle": "string|null",
  "employeeRole": "string|null"
}

EmployeeResearchActivityResponse
{
  "id": 0,
  "researchActivityTitle": "string",
  "employeeRole": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

CreateEmployeeScientificAwardRequest
{
  "employeeID": 0,
  "scientificAwardTitle": "string",
  "givenBy": "string"
}

UpdateEmployeeScientificAwardRequest
{
  "id": 0,
  "scientificAwardTitle": "string|null",
  "givenBy": "string|null"
}

EmployeeScientificAwardResponse
{
  "id": 0,
  "scientificAwardTitle": "string",
  "givenBy": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

CreateEmployeeSocialRequest
{
  "employeeID": 0,
  "socialName": "string",
  "linkToSocial": "string"
}

UpdateEmployeeSocialRequest
{
  "id": 0,
  "socialName": "string|null",
  "linkToSocial": "string|null"
}

EmployeeSocialResponse
{
  "id": 0,
  "socialName": "string",
  "linkToSocial": "string",
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}

CreateEmployeeWorkExperienceRequest
{
  "employeeID": 0,
  "workplace": "string",
  "description": "string",
  "jobTitle": "string",
  "dateStart": "RFC3339",
  "dateEnd": "RFC3339|null",
  "ongoing": true
}

UpdateEmployeeWorkExperienceRequest
{
  "id": 0,
  "jobTitle": "string|null",
  "workplace": "string|null",
  "description": "string|null",
  "dateStart": "RFC3339|null",
  "dateEnd": "RFC3339|null",
  "ongoing": true
}

EmployeeWorkExperienceResponse
{
  "id": 0,
  "workplace": "string",
  "description": "string",
  "jobTitle": "string",
  "dateStart": "RFC3339",
  "dateEnd": "RFC3339",
  "ongoing": true,
  "createdAt": "RFC3339",
  "updatedAt": "RFC3339"
}
