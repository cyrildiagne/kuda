let uiContainer;

const cli = new URL(window.location.href).searchParams.get("cli");

// Setup Firebase UI.
const uiConfig = {
  callbacks: {
    signInSuccessWithAuthResult: (authResult, redirectUrl) => false
  },
  signInOptions: [
    firebase.auth.EmailAuthProvider.PROVIDER_ID,
    firebase.auth.GithubAuthProvider.PROVIDER_ID
  ],
  signInFlow: "popup",
  tosUrl: config.termsOfServiceURL,
  privacyPolicyUrl: config.privacyPolicyURL
};

async function handleAuthChanged(user) {
  if (user) {
    uiContainer.style.display = "none";
    accessToken = await user.getIdToken();
    const userData = JSON.stringify(user, null, "  ");
    // Retrieve CLI argument.
    if (cli) {
      const url = "http://localhost:" + cli;
      try {
        await fetch(url, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: userData
        });
      } catch (e) {
        document.getElementById("account-details").textContent =
          "Error: Authentication successful but could not reach CLI.";
        console.error(e);
        return;
      }
      const msg = "Authentication successful! You can now close this window.";
      document.getElementById("account-details").textContent = msg;
    } else {
      document.getElementById("sign-out").style.display = "inline";
      document.getElementById("account-details").textContent = userData;
    }
  } else {
    // User is signed out.
    document.getElementById("sign-out").style.display = "none";
    document.getElementById("account-details").textContent = "";
    // show Firebase UI.
    uiContainer.style.display = "block";
  }
}

// window.onload = () => {
uiContainer = document.getElementById("firebaseui-auth-container");

firebase.initializeApp({
  apiKey: config.apiKey,
  authDomain: config.authDomain
});

// Listen to change in auth state so it displays the correct UI for when
// the user is signed in or not.
firebase.auth().onAuthStateChanged(handleAuthChanged);

// Signout
const signOutButton = document.getElementById("sign-out");
signOutButton.addEventListener("click", () => {
  firebase
    .auth()
    .signOut()
    .then(res => {
      ui.start("#firebaseui-auth-container", uiConfig);
    });
});

// Start Firebase UI.
ui = new firebaseui.auth.AuthUI(firebase.auth());
ui.start("#firebaseui-auth-container", uiConfig);
if (cli) {
  ui.disableAutoSignIn();
}
