function addFavorite(beatmapsetId) {
    var url = "/users/" + currentUser + "/favourites";

    performApiRequest("POST", url, { set_id: beatmapsetId }, function (xhr) {
        var favourites = document.getElementById("favourites-button");
        favourites.innerHTML = "(Remove Favourite)";
        favourites.style.color = "red";
        favourites.onclick = function () {
            removeFavorite(beatmapsetId);
        };
    });
}

function removeFavorite(beatmapsetId) {
    var url = "/users/" + currentUser + "/favourites/" + beatmapsetId;

    performApiRequest("DELETE", url, null, function (xhr) {
        var favourites = document.getElementById("favourites-button");
        favourites.innerHTML = "(Add Favourite)";
        favourites.style.color = "green";
        favourites.onclick = function () {
            addFavorite(beatmapsetId);
        };
    });
}

function copySetId(element) {
    var set_id = element.getAttribute("setid");

    navigator.clipboard.writeText(set_id).then(
        function () {
            element.innerHTML = "Copied!";
            element.style.color = "green";
        },
        function () {
            element.innerHTML = "Failed to copy!";
            element.style.color = "red";
        }
    );

    setTimeout(function () {
        element.innerHTML = "Copy Beatmapset ID";
        element.style.color = "rgb(0, 102, 204)";
    }, 1500);
}

function updateBeatmapsetMetadata(event) {
    event.preventDefault();

    var data = convertFormToJson(event.target);
    var url = "/beatmapsets/" + data.beatmapset_id;

    performApiRequest(
        "PATCH",
        url,
        data,
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            alert(response.details);
        }
    );
}

function editBeatmapDescription() {
    // TODO: Move this endpoint to the new API
    var description = document.querySelector(".beatmap-description .bbcode");
    if (!description) {
        return;
    }

    var form = document.createElement("form");
    var textarea = document.createElement("textarea");
    textarea.className = "description-editor bbcode-editor";
    textarea.innerHTML = bbcodeDescription;
    textarea.name = "description";
    form.appendChild(textarea);

    var submitButton = document.createElement("input");
    submitButton.type = "submit";
    submitButton.value = "Save";
    form.appendChild(submitButton);

    form.onsubmit = function (event) {
        event.preventDefault();
        var url = "/users/" + currentUser + "/beatmapsets/" + beatmapsetId + "/description";

        performApiRequest(
            "PATCH",
            url,
            { bbcode: textarea.value },
            function (xhr) {
                reloadPageSoon();
            },
            function (xhr) {
                showError(xhr, "An error occurred while trying to update the description.");
            }
        );
    };

    description.replaceWith(form);
}

function convertBanchoSpoilerBoxes() {
    // osu.ppy.sh spoilerbox conversion
    var spoilerBoxes = document.querySelectorAll(".bbcode-spoilerbox");
    for (var i = 0; i < spoilerBoxes.length; i++) {
        $(spoilerBoxes[i]).addClass("spoiler");
    }

    var spoilerBoxContents = document.querySelectorAll(".bbcode-spoilerbox__body");
    for (var i = 0; i < spoilerBoxContents.length; i++) {
        $(spoilerBoxContents[i]).addClass("spoiler-body");
    }

    var spoilerBoxHeads = document.querySelectorAll(".bbcode-spoilerbox__link");
    for (var i = 0; i < spoilerBoxHeads.length; i++) {
        var spoilerBox = spoilerBoxHeads[i];

        // Change element type to div (in older versions of Chrome, use workarounds if necessary)
        var newSpoilerBox = document.createElement("div");
        newSpoilerBox.className = spoilerBox.className + " spoiler-head";
        newSpoilerBox.innerHTML = spoilerBox.innerHTML;
        spoilerBox.parentNode.replaceChild(newSpoilerBox, spoilerBox);

        newSpoilerBox.onclick = function () {
            toggleSpoiler(this);
        };
    }
}

function setBeatmapVolume(volume) {
    var beatmapPreview = document.getElementById("beatmap-preview");
    if (beatmapPreview) {
        beatmapPreview.volume = volume;
    }
}

function acceptCollaborationRequest(beatmapId, requestId) {
    var url = "/beatmaps/" + beatmapId + "/collaborations/requests/" + requestId + "/accept";

    performApiRequest(
        "POST",
        url,
        null,
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            showError(xhr, "An error occurred while trying to accept the invite.");
        }
    );
}

function rejectCollaborationRequest(beatmapId, requestId) {
    if (!confirm("Are you sure you want to decline?")) return;

    var url = "/beatmaps/" + beatmapId + "/collaborations/requests/" + requestId;

    performApiRequest(
        "DELETE",
        url,
        null,
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            showError(xhr, "An error occurred while trying to decline the invite.");
        }
    );
}

function removeCollaborationRequest(beatmapId, requestId) {
    if (!confirm("Are you sure?")) return;

    var url = "/beatmaps/" + beatmapId + "/collaborations/requests/" + requestId;

    performApiRequest(
        "DELETE",
        url,
        null,
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            showError(xhr, "An error occurred while trying to delete your invite.");
        }
    );
}

function createCollaborationRequest(beatmapId) {
    var username = prompt("Enter the username of the user you want to collaborate with:");
    if (!username) {
        return;
    }

    var url = "/users/lookup/" + username.trim();

    performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            user = JSON.parse(xhr.responseText);
            createCollaborationRequestFromUserId(user.id, beatmapId);
        },
        function (xhr) {
            if (xhr.status === 404) {
                alert("User not found. Please check the username and try again.");
            } else {
                alert("An error occurred while looking up the user.");
            }
        }
    );
}

function createCollaborationRequestFromUserId(userId, beatmapId) {
    var url = "/beatmaps/" + beatmapId + "/collaborations/requests";

    performApiRequest(
        "POST",
        url,
        { user_id: userId },
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            showError(xhr, "An error occurred while creating your invite.");
        }
    );
}

function editCollaborationRequest(beatmapId, collaborationId, edits) {
    var url = "/beatmaps/" + beatmapId + "/collaborations/" + collaborationId;

    performApiRequest(
        "PATCH",
        url,
        edits,
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            showError(xhr, "An error occurred while editing the collaboration invite.");
        }
    );
}

function removeCollaboration(beatmapId, collaborationId) {
    if (!confirm("Are you sure you want to remove this collaborator?")) {
        return;
    }

    var url = "/beatmaps/" + beatmapId + "/collaborations/" + collaborationId;

    performApiRequest(
        "DELETE",
        url,
        null,
        function (xhr) {
            reloadPageSoon();
        },
        function (xhr) {
            showError(xhr, "An error occurred while trying to remove the collaborator.");
        }
    );
}

function collaborationMakeAuthor(beatmapId, collaborationId) {
    return editCollaborationRequest(beatmapId, collaborationId, {
        allow_resource_updates: canUpdateResources(collaborationId),
        is_beatmap_author: true
    });
}

function collaborationRemoveAuthor(beatmapId, collaborationId) {
    return editCollaborationRequest(beatmapId, collaborationId, {
        allow_resource_updates: canUpdateResources(collaborationId),
        is_beatmap_author: false
    });
}

function collaborationAllowResourceUpdates(beatmapId, collaborationId) {
    return editCollaborationRequest(beatmapId, collaborationId, {
        allow_resource_updates: true,
        is_beatmap_author: isBeatmapAuthor(collaborationId)
    });
}

function collaborationDisallowResourceUpdates(beatmapId, collaborationId) {
    return editCollaborationRequest(beatmapId, collaborationId, {
        allow_resource_updates: false,
        is_beatmap_author: isBeatmapAuthor(collaborationId)
    });
}

function isBeatmapAuthor(collaborationId) {
    var collaborator = document.getElementById("collaborator-" + collaborationId);
    return collaborator.innerHTML.includes("Remove Author Status");
}

function canUpdateResources(collaborationId) {
    var collaborator = document.getElementById("collaborator-" + collaborationId);
    return collaborator.innerHTML.includes("Disallow Resource Updates");
}

function showError(xhr, defaultMessage) {
    var errorMessage = defaultMessage || "An error occurred while processing your request.";
    try {
        var response = JSON.parse(xhr.responseText);
        errorMessage = response.details;
    } catch (e) {
        console.error("Failed to parse error response:", e);
    }
    alert(errorMessage);
}

var scores = document.querySelectorAll(".scores tbody tr");
for (var i = 0; i < scores.length; i++) {
    $(scores[i]).on("click", function (e) {
        if ($(e.target).closest("a").length) return;

        window.location.href = "/scores/" + this.id;
    });
}

$(document).ready(function () {
    var url = window.location.pathname;
    if (!url.startsWith("/b/") && !url.startsWith("/s/")) {
        return;
    }

    convertBanchoSpoilerBoxes();
    setBeatmapVolume(0.5);

    if (!isBeatmapsetOwner) {
        return;
    }

    var descriptionElements = document.querySelectorAll(".beatmap-description, .beatmap-description *");

    for (var i = 0; i < descriptionElements.length; i++) {
        $(descriptionElements[i]).on("dblclick", function (event) {
            editBeatmapDescription();
        });
    }
});
