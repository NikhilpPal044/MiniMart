const API_URL = "http://localhost:8080/api";
const productList = document.getElementById("productList");
const appMessage = document.getElementById("message");
async function apiRequest(endpoint, options = {}) {
    const token = localStorage.getItem("token");

    const headers = {
        "Content-Type": "application/json",
        ...options.headers
    };

    if (token) {
        headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_URL}${endpoint}`, {
        ...options,
        credentials: "include",
        headers
    });

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
        throw new Error(
            data.message || data.error || `Request failed: ${response.status}`
        );
    }

    return data;


}


async function loadProducts() {
    try {
        const response = await fetch(`${API_URL}/products`);

        if (!response.ok) {
            throw new Error("Could not load products.");
        }

        const data = await response.json();
        const products = data.products;

        productList.innerHTML = "";

        if (!Array.isArray(products) || products.length === 0) {
            productList.innerHTML = "<p>No products available.</p>";
            return;
        }
        products.forEach((product) => {
            const card = document.createElement("article");
            card.className = "product-card";

            // Product image
            const image = document.createElement("img");
            if (product.name === "Wireless Mouse") {
                image.src = "https://images.unsplash.com/photo-1527814050087-3793815479db?w=400";
            } else if (product.name === "Mechanical Keyboard") {
                image.src = "https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=400";
            } else if (product.name === "USB Headphones") {
                image.src = "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400";
            } else {
                image.src = "https://placehold.co/300x200/png?text=Product";
            }


            image.alt = product.name;
            image.style.cssText =
                "display:block;width:100%;height:200px;object-fit:contain;";

            image.onload = function () {
                console.log("Image successfully loaded:", product.image_url);
            };

            image.onerror = function () {
                console.log("Image failed to load:", product.image_url);
            };

            // Product information
            const info = document.createElement("div");
            info.className = "product-info";

            const name = document.createElement("h3");
            name.textContent = product.name;

            const description = document.createElement("p");
            description.textContent = product.description;

            const price = document.createElement("p");
            price.textContent = `Price: ₹${product.price}`;

            const stock = document.createElement("p");
            stock.textContent = `Stock: ${product.stock}`;

            const addButton = document.createElement("button");
            addButton.textContent = "Add to Cart";
            addButton.className = "add-to-cart-btn";

            addButton.addEventListener("click", function () {
                addToCart(product.id);
            });

            info.append(name, description, price, stock, addButton);

            // Add image and information to card
            card.append(image, info);

            productList.appendChild(card);

        });
    } catch (error) {
        productList.innerHTML = "";
        const errorText = document.createElement("p");
        errorText.textContent =
            "Unable to display products. Please check the browser console for details.";
        productList.appendChild(errorText);

        console.error("Product loading error:", error);
    }
}

async function addToCart(productId) {
    console.log("ADD TO CART BUTTON CLICKED:", productId);
    const token = localStorage.getItem("token");
    if (!token) {
        appMessage.textContent = "Please log in first to add items to your cart.";
        return;
    }

    try {
        const response = await fetch(`${API_URL}/cart`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({
                product_id: productId,
                quantity: 1
            })
        });

        const data = await response.json().catch(() => ({}));

        console.log("Add to cart response:", response.status, data);

        if (!response.ok) {
            throw new Error(data.error || data.message || "Could not add product.");
        }

        appMessage.textContent = "Product added to cart!";
    } catch (error) {
        appMessage.textContent = error.message;
        console.error("Add to cart error:", error);
    }

}



async function logoutUser() {
    try {
        await apiRequest("/auth/logout", {
            method: "POST"
        });
    } catch (error) {
        console.error("Logout error:", error.message);
    }

    localStorage.removeItem("user");
    window.location.href = "login.html";
}
if (productList) {
    loadProducts();
}


