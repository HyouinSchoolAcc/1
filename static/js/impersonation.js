(function () {
  "use strict";

  var BANNER_ID = "impersonation-banner";

  function injectBanner(username) {
    if (document.getElementById(BANNER_ID)) return;

    var banner = document.createElement("div");
    banner.id = BANNER_ID;
    Object.assign(banner.style, {
      position: "fixed",
      top: "0",
      left: "0",
      right: "0",
      zIndex: "999999",
      background: "#d63384",
      color: "#fff",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      gap: "12px",
      padding: "6px 16px",
      fontSize: "14px",
      fontFamily: "system-ui, sans-serif",
      boxShadow: "0 2px 8px rgba(0,0,0,0.25)",
    });

    var label = document.createElement("span");
    label.textContent = "Viewing as: " + username;
    label.style.fontWeight = "600";

    var btn = document.createElement("button");
    btn.textContent = "Stop Impersonating";
    Object.assign(btn.style, {
      background: "#fff",
      color: "#d63384",
      border: "none",
      borderRadius: "4px",
      padding: "3px 12px",
      cursor: "pointer",
      fontWeight: "600",
      fontSize: "13px",
    });
    btn.onmouseenter = function () { btn.style.background = "#f0f0f0"; };
    btn.onmouseleave = function () { btn.style.background = "#fff"; };
    btn.onclick = function () {
      fetch("/api/stop-impersonate", { method: "POST" })
        .then(function (res) { return res.json(); })
        .then(function () { location.reload(); });
    };

    banner.appendChild(label);
    banner.appendChild(btn);
    document.body.prepend(banner);

    document.body.style.marginTop = banner.offsetHeight + "px";
  }

  function injectEditorWidget(currentUserId) {
    var wrapper = document.createElement("div");
    wrapper.id = "impersonate-widget";
    Object.assign(wrapper.style, {
      position: "fixed",
      bottom: "16px",
      right: "16px",
      zIndex: "999998",
    });

    var toggle = document.createElement("button");
    toggle.textContent = "Impersonate User";
    Object.assign(toggle.style, {
      background: "#6f42c1",
      color: "#fff",
      border: "none",
      borderRadius: "8px",
      padding: "8px 16px",
      cursor: "pointer",
      fontWeight: "600",
      fontSize: "13px",
      boxShadow: "0 2px 8px rgba(0,0,0,0.2)",
    });

    var panel = document.createElement("div");
    panel.style.display = "none";
    Object.assign(panel.style, {
      position: "absolute",
      bottom: "44px",
      right: "0",
      width: "300px",
      maxHeight: "400px",
      background: "#fff",
      border: "1px solid #dee2e6",
      borderRadius: "8px",
      boxShadow: "0 4px 16px rgba(0,0,0,0.15)",
      overflow: "hidden",
      display: "none",
    });

    var search = document.createElement("input");
    search.type = "text";
    search.placeholder = "Search users...";
    Object.assign(search.style, {
      width: "100%",
      padding: "8px 12px",
      border: "none",
      borderBottom: "1px solid #dee2e6",
      outline: "none",
      fontSize: "13px",
      boxSizing: "border-box",
    });

    var list = document.createElement("div");
    Object.assign(list.style, {
      maxHeight: "340px",
      overflowY: "auto",
    });

    panel.appendChild(search);
    panel.appendChild(list);
    wrapper.appendChild(panel);
    wrapper.appendChild(toggle);
    document.body.appendChild(wrapper);

    var users = [];
    var panelOpen = false;

    toggle.onclick = function () {
      panelOpen = !panelOpen;
      panel.style.display = panelOpen ? "block" : "none";
      if (panelOpen && users.length === 0) loadUsers();
      if (panelOpen) search.focus();
    };

    function loadUsers() {
      list.innerHTML = '<div style="padding:12px;color:#999;text-align:center;">Loading...</div>';
      fetch("/api/users")
        .then(function (r) { return r.json(); })
        .then(function (data) {
          users = data;
          renderList("");
        })
        .catch(function () {
          list.innerHTML = '<div style="padding:12px;color:#dc3545;">Failed to load users</div>';
        });
    }

    function renderList(query) {
      var q = query.toLowerCase();
      var filtered = users.filter(function (u) {
        return u.username.toLowerCase().includes(q) || u.email.toLowerCase().includes(q) || u.id.includes(q);
      });

      list.innerHTML = "";
      if (filtered.length === 0) {
        list.innerHTML = '<div style="padding:12px;color:#999;text-align:center;">No users found</div>';
        return;
      }

      filtered.forEach(function (u) {
        var row = document.createElement("div");
        Object.assign(row.style, {
          padding: "8px 12px",
          cursor: "pointer",
          borderBottom: "1px solid #f0f0f0",
          fontSize: "13px",
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
        });
        if (u.id === currentUserId) {
          row.style.background = "#e8f5e9";
          row.style.fontWeight = "600";
        }
        row.onmouseenter = function () { if (u.id !== currentUserId) row.style.background = "#f8f9fa"; };
        row.onmouseleave = function () { if (u.id !== currentUserId) row.style.background = ""; };

        var info = document.createElement("div");
        info.innerHTML =
          '<div style="font-weight:500;">' + escHtml(u.username) + "</div>" +
          '<div style="font-size:11px;color:#999;">' + escHtml(u.role) + " &middot; " + escHtml(u.email) + "</div>";

        row.appendChild(info);

        if (u.id !== currentUserId) {
          row.onclick = function () { doImpersonate(u.id); };
        }

        list.appendChild(row);
      });
    }

    search.oninput = function () { renderList(search.value); };

    function doImpersonate(userId) {
      fetch("/api/impersonate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user_id: userId }),
      })
        .then(function (r) { return r.json(); })
        .then(function (data) {
          if (data.success) location.reload();
          else alert("Failed: " + (data.error || "unknown error"));
        });
    }
  }

  function escHtml(s) {
    var d = document.createElement("div");
    d.textContent = s;
    return d.innerHTML;
  }

  function init() {
    fetch("/api/current_user")
      .then(function (r) { return r.json(); })
      .then(function (data) {
        if (!data.logged_in) return;
        if (data.impersonating) {
          injectBanner(data.impersonating);
          injectEditorWidget(data.user_id);
        } else if (data.role === "editor") {
          injectEditorWidget(data.user_id);
        }
      })
      .catch(function () { /* silent */ });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
