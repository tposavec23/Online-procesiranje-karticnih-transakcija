package hr.foi.air.sevenpay.core.network

import hr.foi.air.sevenpay.core.BuildConfig
import hr.foi.air.sevenpay.core.data.SessionManager
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

class ApiClient(
    private val sessionManager: SessionManager
) {

    private val okHttpClient = OkHttpClient.Builder()
        .addInterceptor { chain ->
            val token = sessionManager.token
            val request = if (token != null) {
                chain.request().newBuilder()
                    .header("Authorization", "Bearer $token")
                    .build()
            } else {
                chain.request()
            }
            val response = chain.proceed(request)
            if (response.code() == 401 && token != null) {
                sessionManager.expire()
            }
            response
        }
        .build()

    private fun retrofit(baseUrl: String): Retrofit =
        Retrofit.Builder()
            .baseUrl(baseUrl)
            .client(okHttpClient)
            .addConverterFactory(GsonConverterFactory.create())
            .build()

    val catalogueApi: CatalogueApi by lazy { retrofit(BuildConfig.CATALOGUE_BASE_URL).create(CatalogueApi::class.java) }
    val accountApi: AccountApi by lazy { retrofit(BuildConfig.ACCOUNT_BASE_URL).create(AccountApi::class.java) }
}