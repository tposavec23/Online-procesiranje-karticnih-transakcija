package hr.foi.air.sevenpay.core.network

import hr.foi.air.sevenpay.core.BuildConfig
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

class ApiClient {

    private val okHttpClient = OkHttpClient()

    private fun retrofit(baseUrl: String): Retrofit =
        Retrofit.Builder()
            .baseUrl(baseUrl)
            .client(okHttpClient)
            .addConverterFactory(GsonConverterFactory.create())
            .build()

    val catalogApi: CatalogApi by lazy { retrofit(BuildConfig.CATALOG_BASE_URL).create(CatalogApi::class.java) }
    val accountApi: AccountApi by lazy { retrofit(BuildConfig.ACCOUNT_BASE_URL).create(AccountApi::class.java) }
}
