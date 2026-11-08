import java.io.StringReader
import java.util.Properties

plugins {
    alias(libs.plugins.android.library)
}

// Addresses of the backend services the app connects to.
// You can set them in android/local.properties (not in git), each key on its own line,
// for example to run the app on a phone on the same network as the computer:
//   CATALOG_BASE_URL=http://192.168.1.10:8081/
//   ACCOUNT_BASE_URL=http://192.168.1.10:8082/
// Without these keys the defaults below are used (10.0.2.2 is this computer as seen from the emulator).
// Every address must end with "/".
val localProperties = providers.fileContents(
    rootProject.layout.projectDirectory.file("local.properties")
).asText.map { text -> Properties().apply { load(StringReader(text)) } }

fun backendUrl(key: String, default: String): String {
    val url = localProperties.map { it.getProperty(key) }.orElse(default).get()
    require(url.endsWith("/")) { "$key in local.properties must end with '/': '$url'" }
    return url
}

android {
    namespace = "hr.foi.air.sevenpay.core"
    compileSdk {
        version = release(37)
    }

    defaultConfig {
        minSdk = 24

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"

        buildConfigField("String", "CATALOG_BASE_URL", "\"${backendUrl("CATALOG_BASE_URL", "http://10.0.2.2:8081/")}\"")
        buildConfigField("String", "ACCOUNT_BASE_URL", "\"${backendUrl("ACCOUNT_BASE_URL", "http://10.0.2.2:8082/")}\"")
    }
    // Debug and release builds currently use the same addresses.
    // Once the 7Pay production server is set up, give the release build its addresses here
    // (buildTypes { release { buildConfigField(...) } }).
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
    buildFeatures {
        buildConfig = true
    }
}

dependencies {
    implementation(libs.retrofit)
    implementation(libs.retrofit.converter.gson)
    implementation(libs.androidx.appcompat)
    implementation(libs.androidx.core.ktx)
    implementation(libs.material)
    testImplementation(libs.junit)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(libs.androidx.junit)
}
