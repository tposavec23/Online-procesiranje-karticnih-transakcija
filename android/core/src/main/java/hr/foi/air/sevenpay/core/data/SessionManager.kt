package hr.foi.air.sevenpay.core.data

class SessionManager {

    @Volatile
    var token: String? = null
        private set

    @Volatile
    var onSessionExpired: (() -> Unit)? = null

    fun saveToken(token: String) {
        this.token = token
    }

    fun clear() {
        token = null
    }

    fun expire() {
        token = null
        onSessionExpired?.invoke()
    }
}