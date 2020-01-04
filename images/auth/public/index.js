let uiContainer;
let currentUser;

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

function handleAuthChanged(user) {
  if (user) {
    currentUser = user;
    uiContainer.style.display = "none";
    user.getIdToken().then(accessToken => {
      document.getElementById("sign-out").style.display = "inline";
      document.getElementById("account-details").textContent = JSON.stringify(
        currentUser,
        null,
        "  "
      );
    });
  } else {
    currentUser = null;
    // User is signed out.
    document.getElementById("sign-out").style.display = "none";
    document.getElementById("account-details").textContent = "";
    // show Firebase UI.
    uiContainer.style.display = "block";
  }
}

window.onload = () => {
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
};
