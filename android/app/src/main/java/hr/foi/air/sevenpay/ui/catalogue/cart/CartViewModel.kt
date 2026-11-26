package hr.foi.air.sevenpay.ui.cart

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.ViewModelProvider.AndroidViewModelFactory.Companion.APPLICATION_KEY
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import hr.foi.air.sevenpay.SevenPayApplication
import hr.foi.air.sevenpay.core.data.CartRepository
import hr.foi.air.sevenpay.core.exception.ArticleNotFoundException
import hr.foi.air.sevenpay.core.exception.OutOfStockException
import hr.foi.air.sevenpay.core.model.Article
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.receiveAsFlow
import kotlinx.coroutines.launch

class CartViewModel(
    private val userId: Int,
    private val cartRepository: CartRepository
) : ViewModel() {


    private val _messages = Channel<String>(Channel.BUFFERED)
    val messages: Flow<String> = _messages.receiveAsFlow()

    private var isAdding = false

    fun addToCart(article: Article, quantity: Int) {
        if (isAdding) return
        isAdding = true

        viewModelScope.launch {
            val message = try {
                cartRepository.addToCart(userId, article.id, quantity)
                if (quantity > 1) {
                    "Dodano u košaricu: ${article.name} × $quantity"
                } else {
                    "Dodano u košaricu: ${article.name}"
                }
            } catch (e: ArticleNotFoundException) {
                "Artikl više ne postoji."
            } catch (e: OutOfStockException) {
                "Nema dovoljno artikala na zalihi."
            } catch (e: CancellationException) {
                throw e
            } catch (e: Exception) {
                "Dodavanje u košaricu nije uspjelo. Provjerite vezu i pokušajte ponovno."
            } finally {
                isAdding = false
            }
            _messages.send(message)
        }
    }

    companion object {
        fun factory(userId: Int): ViewModelProvider.Factory = viewModelFactory {
            initializer {
                val application = this[APPLICATION_KEY] as SevenPayApplication
                CartViewModel(userId, application.container.cartRepository)
            }
        }
    }
}