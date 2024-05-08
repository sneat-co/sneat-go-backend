# Kills processes that occupy ports used by Sneat.app local dev instance
#   4300 - Backend API @ local Google App Engine emulator
#   8070 - Firebase Emulator UI
#   8080 - Firestore Admin UI
#   9099 - Firebase Authentication
kill -9 "$(lsof -ti:4300,8070,8080,9099)".
