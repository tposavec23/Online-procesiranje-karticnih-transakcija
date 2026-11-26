package hr.foi.air.sevenpay

import hr.foi.air.sevenpay.core.data.AccountRepository
import hr.foi.air.sevenpay.core.data.CartRepository
import hr.foi.air.sevenpay.core.data.CatalogueRepository
import hr.foi.air.sevenpay.core.data.NetworkAccountRepository
import hr.foi.air.sevenpay.core.data.NetworkCartRepository
import hr.foi.air.sevenpay.core.data.NetworkCatalogueRepository
import hr.foi.air.sevenpay.core.data.SessionManager
import hr.foi.air.sevenpay.core.network.ApiClient

class AppContainer {
    val sessionManager = SessionManager()
    private val apiClient = ApiClient(sessionManager)

    val accountRepository: AccountRepository by lazy { NetworkAccountRepository(apiClient.accountApi, sessionManager) }
    val cartRepository: CartRepository by lazy { NetworkCartRepository(apiClient.cartApi) }
    val catalogueRepository: CatalogueRepository by lazy { NetworkCatalogueRepository(apiClient.catalogueApi) }
}