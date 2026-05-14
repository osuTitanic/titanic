function fetchIrcToken(onSuccess, onFailure) {
    if (!isLoggedIn()) {
        console.error("User is not logged in, cannot fetch IRC token.");
        return;
    }

    return performApiRequest("GET", "/account/irc/token", null, onSuccess, onFailure);
}

function regenerateIrcToken(onSuccess, onFailure) {
    if (!isLoggedIn()) {
        console.error("User is not logged in, cannot fetch IRC token.");
        return;
    }

    return performApiRequest("POST", "/account/irc/token", null, onSuccess, onFailure);
}

function showIrcTokenInSettings(action) {
    var tokenResolver = action == "view" ? fetchIrcToken : regenerateIrcToken;

    tokenResolver(
        function (xhr) {
            var tokenContainer = document.getElementById("irc-token-container");
            if (!tokenContainer) {
                console.error("IRC token container not found in the document.");
                return;
            }

            var response = JSON.parse(xhr.responseText);
            if (!response || !response.token) {
                console.error("Invalid response format or missing token:", response);
                alert("Failed to retrieve your IRC password. Please try again later!");
                return;
            }

            tokenContainer.innerText = response.token;
        },
        function (xhr) {
            console.error("Failed to fetch IRC token:", xhr);
            alert("Failed to retrieve your IRC password. Please try again later!");
        }
    );
}
