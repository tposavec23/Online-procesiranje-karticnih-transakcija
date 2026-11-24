package hr.foi.air.sevenpay.core.data

class SessionManager {

    @Volatile
    var token: String? = null
        private set

    fun saveToken(token: String) {
        this.token = token
    }

    fun clear() {
        token = null
    }
}