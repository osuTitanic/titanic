function handleNotificationClick(event, linkElement) {
    event.preventDefault();

    var notificationElement = linkElement.closest(".notification");
    var notificationId = notificationElement.id;
    var targetUrl = linkElement.href;

    // Mark notification as read
    performApiRequest("DELETE", "/account/notifications/" + notificationId, null, function (xhr) {
        notificationElement.classList.remove("new");

        // Redirect to the target URL after marking as read
        if (targetUrl) window.location.href = targetUrl;
    });

    return false;
}

function handleNotificationBoxClick(event, notificationId) {
    event.preventDefault();

    var notificationElement = document.getElementById(notificationId);
    var targetUrl = notificationElement.querySelector(".notification-link").href;

    // Mark notification as read
    performApiRequest("DELETE", "/account/notifications/" + notificationId, null, function (xhr) {
        notificationElement.classList.remove("new");

        // Redirect to the target URL after marking as read
        if (targetUrl) window.location.href = targetUrl;
    });

    return false;
}

function clearAllNotifications() {
    performApiRequest("DELETE", "/account/notifications", null, function (xhr) {
        var elements = document.querySelectorAll(".notification");
        for (var i = 0; i < elements.length; i++) {
            elements[i].classList.remove("new");
            elements[i].onclick = function () {};
        }

        reloadPageSoon(500);
    });
}

function removeBookmark(topicId) {
    performApiRequest("DELETE", "/forum/bookmarks/" + topicId, null, function (xhr) {
        var data = JSON.parse(xhr.responseText);
        var bookmark = document.getElementById("bookmark-" + topicId);
        bookmark.parentNode.removeChild(bookmark);

        var totalBookmarks = document.querySelectorAll(".bookmark").length;
        if (totalBookmarks === 0) {
            var bookmarks = document.querySelector(".bookmarks");
            bookmarks.innerHTML = "You have no bookmarks";
        }
    });
    return false;
}

function removeFriend(element) {
    if (!isLoggedIn()) return;

    performApiRequest("DELETE", "/account/friends?id=" + element.id, null, function (xhr) {
        var data = JSON.parse(xhr.responseText);
        var friendContainer = element.parentElement.parentElement;
        friendContainer.style.opacity = 0;

        setTimeout(function () {
            friendContainer.parentNode.removeChild(friendContainer);
        }, 350);
    });
    return false;
}

function dataExport() {
    var passwordValidation = document.getElementById("data-export-password");
    var password = passwordValidation.value;
    if (!password) {
        alert("Please enter your current password to confirm the data export.");
        return;
    }
    var data = { password: password };

    var recaptchaResponse = document.getElementById("recaptcha-response");
    if (recaptchaResponse && recaptchaResponse.value) {
        data["recaptcha_response"] = recaptchaResponse.value;
    }

    var exportButton = document.getElementById("data-export-button");
    exportButton.disabled = true;
    exportButton.value = "Exporting...";

    performApiRequest(
        "POST",
        "/account/export",
        data,
        function (xhr) {
            var blob = new Blob([xhr.responseText], { type: "application/json" });
            var url = URL.createObjectURL(blob);

            // Invoke download prompt with json response data
            var a = document.createElement("a");
            a.download = "data_export.json";
            a.href = url;
            a.click();

            exportButton.value = "Exported!";
        },
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            if (!response.details) {
                alert("An error occurred while exporting your data. Please try again later.");
            } else {
                alert(response.details);
            }

            exportButton.disabled = false;
            exportButton.value = "Export";
        }
    );
}
