package api4userus

//var createUser = facade4userus.CreateUser

// httpPostCreateUser creates user record
//func httpPostCreateUser(w http.ResponseWriter, r *http.Request) {
//	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, endpoints.VerifyRequest{
//		MinContentLength: apicore.MinJSONRequestSize,
//		MaxContentLength: 1 * apicore.KB,
//		AuthRequired:     true,
//	})
//	if err != nil {
//		return
//	}
//	var request facade4userus.CreateUserRequest
//	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
//		return
//	}
//	requestWithLogInfo := facade4userus.CreateUserRequestWithRemoteClientInfo{
//		CreateUserRequest: request,
//		RemoteClient:      apicore.GetRemoteClientInfo(r),
//	}
//	//if cityLatLong := r.Header.Get("X-Appengine-CityLatLong"); cityLatLong != "" {
//	//	latLong := strings.Split(cityLatLong, ",")
//	//	var lat, long float64
//	//	if lat, err = strconv.ParseFloat(latLong[0], 64); err != nil {
//	//		err = validation.NewErrBadRequestFieldValue("request.header[X-Appengine-CityLatLong]", "not valid lat")
//	//		return
//	//	}
//	//	if long, err = strconv.ParseFloat(latLong[0], 64); err != nil {
//	//		err = validation.NewErrBadRequestFieldValue("request.header[X-Appengine-CityLatLong]", "not valid long")
//	//		return
//	//	}
//	//	if !(lat == 0 && long == 0) {
//	//		requestWithLogInfo.GeoCityPoint = &appengine.GeoPoint{
//	//			Lat: lat,
//	//			Lng: long,
//	//		}
//	//	}
//	//}
//	err = createUser(ctx, userContext.ContactID(), requestWithLogInfo)
//	apicore.IfNoErrorReturnCreatedOK(ctx, w, err)
//}
