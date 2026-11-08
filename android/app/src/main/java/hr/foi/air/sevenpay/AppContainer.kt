package hr.foi.air.sevenpay

import hr.foi.air.sevenpay.core.data.AccountRepository
import hr.foi.air.sevenpay.core.data.CatalogRepository
import hr.foi.air.sevenpay.core.data.NetworkAccountRepository
import hr.foi.air.sevenpay.core.data.NetworkCatalogRepository
import hr.foi.air.sevenpay.core.network.ApiClient

class AppContainer {
    private val apiClient = ApiClient()

    val accountRepository: AccountRepository by lazy { NetworkAccountRepository(apiClient.accountApi) }
    val catalogRepository: CatalogRepository by lazy { NetworkCatalogRepository(apiClient.catalogApi) }
}
