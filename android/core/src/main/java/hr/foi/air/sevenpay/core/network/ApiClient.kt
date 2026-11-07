package hr.foi.air.sevenpay.core.network

import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

object ApiClient {

    private const val BASE_URL = "http://10.0.2.2:8081/"
    private const val ACCOUNT_BASE_URL = "http://10.0.2.2:8082/"

    val catalogApi: CatalogApi by lazy {
        Retrofit.Builder()
            .baseUrl(BASE_URL)
            .addConverterFactory(GsonConverterFactory.create())
            .build()
            .create(CatalogApi::class.java)
    }

    val accountApi: AccountApi by lazy {
        Retrofit.Builder()
            .baseUrl(ACCOUNT_BASE_URL)
            .addConverterFactory(GsonConverterFactory.create())
            .build()
            .create(AccountApi::class.java)
    }
}