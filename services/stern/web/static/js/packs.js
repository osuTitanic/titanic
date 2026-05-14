function expandBeatmapPack(id) {
    if (id === null || id === undefined) return;
    var element = document.getElementById("pack-" + id);

    if (element.className.indexOf("expanded") !== -1) {
        $(element).removeClass("expanded");
        $(element).slideUp(500);
        return;
    }

    if (element.className.indexOf("loaded") !== -1) {
        $(element).addClass("expanded");
        $(element).slideDown(500);
        return;
    }

    element.innerHTML = "<center>Loading...</center>";
    $(element).addClass("expanded");
    element.style.height = "";
    $(element).hide().slideDown(100);
    loadBeatmapPackInfo(id, function () {
        $(element).hide().slideDown(500);
    });
}

function loadBeatmapPackInfo(id, callback) {
    var element = document.getElementById("pack-" + id);
    var url = "/beatmapsets/packs/" + currentCategory + "/" + id;

    performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var info = JSON.parse(xhr.responseText);
            var heading = document.createElement("h2");
            heading.innerText = info.name;
            heading.style.marginTop = "10px";

            var creatorLink = document.createElement("a");
            creatorLink.href = "/u/" + info.creator.id;
            creatorLink.innerText = info.creator.name;

            var uploadedDate = new Date(info.created_at);
            var uploadedDateString = uploadedDate.toLocaleDateString("de-DE", {
                year: "numeric",
                month: "numeric",
                day: "numeric",
                hour: "numeric",
                minute: "numeric",
                second: "numeric"
            });

            var uploadedDateInfo = document.createElement("time");
            uploadedDateInfo.setAttribute("datetime", uploadedDateString);
            uploadedDateInfo.innerText = info.created_at;

            var uploadedInfo = document.createElement("p");
            uploadedInfo.appendChild(document.createTextNode("Created by "));
            uploadedInfo.appendChild(creatorLink);
            uploadedInfo.appendChild(document.createTextNode(" on "));
            uploadedInfo.appendChild(uploadedDateInfo);

            var description = document.createElement("p");
            description.innerText = info.description;

            var beatmapList = document.createElement("ul");
            beatmapList.style.marginTop = "10px";

            for (var i = 0; i < info.entries.length; i++) {
                var beatmapset = info.entries[i].beatmapset;
                var beatmapLink = document.createElement("a");
                beatmapLink.href = "/s/" + beatmapset.id;
                beatmapLink.innerText = beatmapset.artist + " - " + beatmapset.title + " (" + beatmapset.creator + ")";
                var beatmapItem = document.createElement("li");
                beatmapItem.appendChild(beatmapLink);
                beatmapList.appendChild(beatmapItem);
            }

            $(element).addClass("loaded");
            element.innerHTML = "";
            element.appendChild(heading);
            element.appendChild(uploadedInfo);
            element.appendChild(description);
            element.appendChild(beatmapList);
            element.appendChild(document.createElement("br"));
            if (callback) callback();
        },
        function (xhr) {
            element.innerHTML = "<center>Failed to load beatmap pack info.</center>";
            console.error("Failed to load beatmap pack info: " + xhr.responseText);
        }
    );
}
