function handleNotificationClick(event, linkElement) {
    event.preventDefault();

    var linkElement = $(linkElement);
    var notificationElement = linkElement.closest(".notification");
    var notificationId = notificationElement.attr("id");
    var targetUrl = linkElement.attr("href");

    // Mark notification as read
    performApiRequest("DELETE", "/account/notifications/" + notificationId, null, function (xhr) {
        notificationElement.removeClass("new");

        // Redirect to the target URL after marking as read
        if (targetUrl) window.location.href = targetUrl;
    });

    return false;
}

function handleNotificationBoxClick(event, notificationId) {
    event.preventDefault();

    var notificationElement = $("#" + notificationId);
    var targetUrl = notificationElement.find(".notification-link").attr("href");

    // Mark notification as read
    performApiRequest("DELETE", "/account/notifications/" + notificationId, null, function (xhr) {
        notificationElement.removeClass("new");

        // Redirect to the target URL after marking as read
        if (targetUrl) window.location.href = targetUrl;
    });

    return false;
}

function clearAllNotifications() {
    performApiRequest("DELETE", "/account/notifications", null, function (xhr) {
        $(".notification").removeClass("new").removeAttr("onclick");
        reloadPageSoon(500);
    });
}

function removeBookmark(topicId) {
    performApiRequest("DELETE", "/forum/bookmarks/" + topicId, null, function (xhr) {
        $("#bookmark-" + topicId).remove();
        if ($(".bookmark").length === 0) {
            $(".bookmarks").text("You have no bookmarks");
        }
    });
    return false;
}

function lookupUser(username, callbackSuccess, callbackError) {
    if (!username) {
        if (callbackError) callbackError(null);
        return;
    }

    performApiRequest(
        "GET",
        "/users/lookup/" + encodeURIComponent(username),
        null,
        function (xhr) {
            var user = $.parseJSON(xhr.responseText);
            if (callbackSuccess) callbackSuccess(user);
        },
        function (xhr) {
            if (callbackError) callbackError(xhr);
        }
    );
}

function lookupUserFromSettings(form, selfMessage, callback) {
    if (!isLoggedIn()) return false;

    var formElement = $(form);
    var username = $.trim(formElement.find("[name=username]").val());
    var button = formElement.find("button").first();

    if (!username) {
        alert("Please enter a username.");
        return false;
    }

    button.prop("disabled", true);
    lookupUser(
        username,
        function (user) {
            if (user.id === currentUser) {
                button.prop("disabled", false);
                alert(selfMessage);
                return;
            }

            callback(user, button);
        },
        function (xhr) {
            button.prop("disabled", false);
            apiErrorAlert(xhr, "The user could not be found.");
        }
    );
    return false;
}

function addFriendFromSettings(form) {
    return lookupUserFromSettings(
        form,
        "You cannot add yourself as a friend.",
        function (user, button) {
            addFriend(
                user.id,
                function () {
                    reloadPageSoon(250);
                },
                function (xhr) {
                    button.prop("disabled", false);
                    apiErrorAlert(xhr, "The user could not be added as a friend.");
                }
            );
        }
    );
}

function removeFriendFromSettings(element) {
    return removeRelationshipFromSettings(
        element,
        removeFriend,
        "The user could not be removed from your friends."
    );
}

function addFoeFromSettings(form) {
    return lookupUserFromSettings(
        form,
        "You cannot block yourself.",
        function (user, button) {
            if (!window.confirm("Block " + user.name + "?")) {
                button.prop("disabled", false);
                return;
            }

            addFoe(
                user.id,
                function () {
                    reloadPageSoon(250);
                },
                function (blockXhr) {
                    button.prop("disabled", false);
                    apiErrorAlert(blockXhr, "The user could not be blocked.");
                }
            );
        }
    );
}

function removeFoeFromSettings(element) {
    return removeRelationshipFromSettings(
        element,
        removeFoe,
        "The user could not be unblocked."
    );
}

function removeRelationshipFromSettings(element, removeRelationship, errorMessage) {
    if (!isLoggedIn()) return false;

    var relationshipElement = $(element);
    return removeRelationship(
        relationshipElement.attr("id"),
        function () {
            var relationshipCard = relationshipElement.closest(".relationship-card");
            relationshipCard.fadeOut(350, function () {
                relationshipCard.remove();
            });
        },
        function (xhr) {
            apiErrorAlert(xhr, errorMessage);
        }
    );
}

function dataExport() {
    var password = $("#data-export-password").val();
    if (!password) {
        alert("Please enter your current password to confirm the data export.");
        return;
    }
    var data = { password: password };

    var recaptchaResponse = $("#recaptcha-response").val();
    if (recaptchaResponse) {
        data["recaptcha_response"] = recaptchaResponse;
    }

    var exportButton = $("#data-export-button");
    exportButton.prop("disabled", true).val("Exporting...");

    performApiRequest(
        "POST",
        "/account/export",
        data,
        function (xhr) {
            var blob = new Blob([xhr.responseText], { type: "application/json" });
            var url = URL.createObjectURL(blob);

            // Invoke download prompt with json response data
            var downloadLink = $("<a>")
                .attr("download", "data_export.json")
                .attr("href", url);
            downloadLink[0].click();
            exportButton.val("Exported!");
        },
        function (xhr) {
            apiErrorAlert(xhr, "An error occurred while exporting your data. Please try again later.");
            exportButton.prop("disabled", false).val("Export");
        }
    );
}
