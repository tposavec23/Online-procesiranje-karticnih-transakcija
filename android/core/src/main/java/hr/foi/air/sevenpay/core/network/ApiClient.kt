package hr.foi.air.sevenpay.core.network

import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

object ApiClient {

    private const val BASE_URL = "http://10.0.2.2:8081/"
    private const val ACCOUNT_BASE_URL = "http://10.0.2.2:8082/"

    private fun retrofit(baseUrl: String): Retrofit =
        Retrofit.Builder()
            .baseUrl(baseUrl)
            .addConverterFactory(GsonConverterFactory.create())
            .build()

    val catalogApi: CatalogApi by lazy { retrofit(BASE_URL).create(CatalogApi::class.java) }
    val accountApi: AccountApi by lazy { retrofit(ACCOUNT_BASE_URL).create(AccountApi::class.java) }
}