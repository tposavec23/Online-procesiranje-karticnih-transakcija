package hr.foi.air.sevenpay.ui.demo

import android.content.Context
import android.content.Intent
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.core.content.FileProvider
import hr.foi.air.sevenpay.core.data.CatalogueRepository
import java.io.File
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.launch

@Composable
fun DemoScreen(repository: CatalogueRepository, modifier: Modifier = Modifier) {
    val context = LocalContext.current
    val scope = rememberCoroutineScope()
    var message by remember { mutableStateOf("") }

    val pickCsv = rememberLauncherForActivityResult(ActivityResultContracts.GetContent()) { uri ->
        if (uri == null) return@rememberLauncherForActivityResult
        scope.launch {
            message = try {
                val csv = context.contentResolver.openInputStream(uri)?.use { it.readBytes() }
                    ?: error("datoteka se ne može pročitati")
                val result = repository.importArticles(csv)
                buildString {
                    append("Uvezeno artikala: ${result.imported}")
                    result.errors.orEmpty().forEach { append("\nredak ${it.row}: ${it.reason}") }
                }
            } catch (e: CancellationException) {
                throw e
            } catch (e: Exception) {
                "Uvoz nije uspio: ${e.message}"
            }
        }
    }

    Column(
        modifier = modifier
            .fillMaxSize()
            .padding(24.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        Text("XP radionica: uvoz i izvoz artikala", fontSize = 22.sp, fontWeight = FontWeight.Bold)
        Button(onClick = { pickCsv.launch("*/*") }, modifier = Modifier.fillMaxWidth()) {
            Text("Uvezi CSV")
        }
        Button(
            onClick = {
                scope.launch {
                    message = try {
                        context.shareCsv(repository.exportArticles())
                        "Izvoz je spreman za slanje."
                    } catch (e: CancellationException) {
                        throw e
                    } catch (e: Exception) {
                        "Izvoz nije uspio: ${e.message}"
                    }
                }
            },
            modifier = Modifier.fillMaxWidth()
        ) {
            Text("Pošalji e-poštom")
        }
        Text(message)
    }
}

private fun Context.shareCsv(csv: String) {
    val dir = File(cacheDir, "exports").apply { mkdirs() }
    val file = File(dir, "artikli.csv").apply { writeText(csv) }
    val uri = FileProvider.getUriForFile(this, "$packageName.fileprovider", file)
    val intent = Intent(Intent.ACTION_SEND).apply {
        type = "text/csv"
        putExtra(Intent.EXTRA_STREAM, uri)
        putExtra(Intent.EXTRA_SUBJECT, "Artikli")
        addFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION)
    }
    startActivity(Intent.createChooser(intent, "Pošalji e-poštom"))
}
