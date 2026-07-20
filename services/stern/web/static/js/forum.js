function addSubscription(topicId) {
    performApiRequest("POST", "/forum/subscriptions", { topic_id: topicId }, function (xhr) {
        var button = document.getElementById("subscribe-button");
        button.onclick = function () {
            removeSubscription(topicId);
        };
        $(button).text("Unsubscribe topic");
    });
    return false;
}

function removeSubscription(topicId) {
    performApiRequest("DELETE", "/forum/subscriptions/" + topicId, null, function (xhr) {
        var button = document.getElementById("subscribe-button");
        button.onclick = function () {
            addSubscription(topicId);
        };
        $(button).text("Subscribe topic");
    });
    return false;
}

function addBookmark(topicId) {
    performApiRequest("POST", "/forum/bookmarks", { topic_id: topicId }, function (xhr) {
        var button = document.getElementById("bookmark-button");
        button.onclick = function () {
            removeBookmark(topicId);
        };
        $(button).text("Remove Bookmark");
    });
    return false;
}

function removeBookmark(topicId) {
    performApiRequest("DELETE", "/forum/bookmarks/" + topicId, null, function (xhr) {
        var button = document.getElementById("bookmark-button");
        button.onclick = function () {
            addBookmark(topicId);
        };
        $(button).text("Bookmark topic");
    });
    return false;
}

function deletePost(postId) {
    if (!confirm("Are you sure you want to delete this post?")) return;

    var url = "/forum/0/topics/0/posts/" + postId;

    performApiRequest(
        "DELETE",
        url,
        null,
        function (xhr) {
            var post = document.getElementById("post-" + postId);
            post.innerHTML = "[ Deleted ]";

            var buttons = $(post).parent().find(".post-buttons");
            if (buttons) {
                for (var i = 0; i < buttons.length; i++) {
                    buttons[i].parentNode.removeChild(buttons[i]);
                }
            }
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to delete post.");
        }
    );
    return false;
}

function giveKudos(postId, beatmapsetId) {
    performApiRequest(
        "POST",
        "/beatmapsets/" + beatmapsetId + "/kudosu/" + postId + "/reward",
        null,
        function (xhr) {
            var data = JSON.parse(xhr.responseText);
            var kudosuStatus = document.getElementById("kudosu-status-" + postId);
            var kudosuActions = kudosuStatus.parentElement.getElementsByTagName("a");

            while (kudosuActions[0]) {
                kudosuActions[0].parentNode.removeChild(kudosuActions[0]);
            }

            $(kudosuStatus).text("Earned " + data.amount + " kudosu.");
            kudosuStatus.style.color = "green";
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to give kudosu.");
        }
    );
    return false;
}

function revokeKudos(postId, beatmapsetId) {
    performApiRequest(
        "POST",
        "/beatmapsets/" + beatmapsetId + "/kudosu/" + postId + "/revoke",
        null,
        function (xhr) {
            var kudosuStatus = document.getElementById("kudosu-status-" + postId);
            var kudosuActions = kudosuStatus.parentElement.getElementsByTagName("a");

            while (kudosuActions[0]) {
                kudosuActions[0].parentNode.removeChild(kudosuActions[0]);
            }

            $(kudosuStatus).text("Successfully revoked kudosu.");
            kudosuStatus.style.color = "red";
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to revoke kudosu.");
        }
    );
    return false;
}

function resetKudos(postId, beatmapsetId) {
    performApiRequest(
        "POST",
        "/beatmapsets/" + beatmapsetId + "/kudosu/" + postId + "/reset",
        null,
        function (xhr) {
            var kudosuStatus = document.getElementById("kudosu-status-" + postId);
            var kudosuActions = kudosuStatus.parentElement.getElementsByTagName("a");

            while (kudosuActions[0]) {
                kudosuActions[0].parentNode.removeChild(kudosuActions[0]);
            }

            $(kudosuStatus).text("Successfully reset kudosu.");
            kudosuStatus.style.color = "blue";
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to reset kudosu.");
        }
    );
    return false;
}

function jumpToPage() {
    var page = prompt("Enter the page to jump to:");

    if (page !== null && !isNaN(page) && page > 0) {
        var query = new URLSearchParams();
        query.set("page", page);
        location.search = query.toString();
    }
}

$(document).on("keydown", function (event) {
    var keyCode = event.which || event.keyCode;

    if (event.ctrlKey && keyCode === 13) {
        var form = $(".quick-reply")[0];
        if (!form) return;

        var textarea = $(form).find("textarea")[0];
        if (!textarea) return;

        // Focus textarea if no message was written
        if (trimString(textarea.value) === "") {
            textarea.focus();
            return;
        }

        // Submit quick reply otherwise
        form.submit();
    }
});

function updateTopic(topicId, updates, successCallback, errorCallback) {
    performApiRequest("PATCH", "/forum/0/topics/" + topicId, updates, successCallback, errorCallback);
}

function moveTopic(topicId, targetForumId) {
    updateTopic(
        topicId,
        { forum_id: targetForumId },
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to move topic.");
        }
    );
    return false;
}

function lockTopic(topicId) {
    updateTopic(
        topicId,
        { lock_topic: true },
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to lock topic.");
        }
    );
    return false;
}

function unlockTopic(topicId) {
    updateTopic(
        topicId,
        { lock_topic: false },
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to unlock topic.");
        }
    );
    return false;
}

function setTopicType(topicId, type) {
    updateTopic(
        topicId,
        { type: type },
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to set topic type.");
        }
    );
    return false;
}

function setTopicStatus(topicId, status) {
    updateTopic(
        topicId,
        { status_text: status },
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to set topic status.");
        }
    );
    return false;
}

function setTopicIcon(topicId, iconId) {
    updateTopic(
        topicId,
        { icon_id: iconId },
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to set topic icon.");
        }
    );
    return false;
}

function setTopicTitle(topicId, title) {
    updateTopic(
        topicId,
        { title: title },
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            apiErrorAlert(xhr, "Failed to set topic title.");
        }
    );
    return false;
}
