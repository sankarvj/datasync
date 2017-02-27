package adapter

// func shapeMaker(db *sql.DB, newShaper Shaper, oldShaper Shaper) ([]byte, error) {
// 	//Preserve the old with the lastid ! Copy. Don't refer
// 	sOld.setId(lastid)
// 	//Set Id with serverid
// 	s.setId(serverid)
// 	s.setSynced(false)
// 	//Set Tripid with serverTripid
// 	if s.getTripId() != 0 {
// 		servertripid, available := serverkey(db, "trips", s.getTripId())
// 		if !available {
// 			return nil, errors.New("trip unavailable")
// 		} else {
// 			s.setTripId(servertripid)
// 		}
// 	}
// 	//Set Assignid with serverMemberid
// 	if s.getAssignId() != 0 {
// 		serverAssignid, available := serverkey(db, "members", s.getAssignId())
// 		if !available {
// 			return nil, errors.New("member unavailable")
// 		} else {
// 			s.setAssignId(serverAssignid)
// 		}
// 	}

// 	//Set Categoryid with serverCategoryid
// 	if s.getCategoryId() != 0 {
// 		serverCategoryid, available := serverkey(db, "categories", s.getCategoryId())
// 		if !available {
// 			return nil, errors.New("category unavailable")
// 		} else {
// 			s.setCategoryId(serverCategoryid)
// 		}
// 	}

// 	jsonbody, err := json.Marshal(s)
// 	return jsonbody, err
// }

// func shapeSolver(db *sql.DB, response Response, sNew Basemodel, sOld Basemodel, tablename string, callback ClientCallback) bool {
// 	if response.Id == ResponseSuccess {
// 		shapeRevert(db, sNew, sOld, tablename)
// 		return true
// 	} else if response.Id == ResponseNetworkError {
// 		callback.OnError(response.Id, response.Msg)
// 		return false
// 	} else {
// 		shapeDelete(db, sOld.getId(), tablename)
// 		callback.OnError(response.Id, response.Msg)
// 		return false
// 	}
// }
