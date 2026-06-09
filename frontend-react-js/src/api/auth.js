import {
  getCurrentUser,
  fetchUserAttributes,
  fetchAuthSession,
  signOut as amplifySignOut,
  signIn as amplifySignIn,
} from "aws-amplify/auth";

// Única fonte da verdade — busca usuário E token juntos
export async function checkAuth(forceRefresh = false) {
  try {
    // 1. Verifica se há usuárfo logado
    await getCurrentUser();

    // 2. Busca atributos e token em paralelo
    const [attributes, session] = await Promise.all([
      fetchUserAttributes(),
      fetchAuthSession({ forceRefresh }),
    ]);

    const idToken = session.tokens?.idToken?.toString();

    if (!idToken) throw new Error("Token não encontrado");



    return {attributes, token: idToken};
  } catch (err) {
    console.error(err);
    return null;
  }
}
export async function signOut() {
  try {
    await amplifySignOut({ global: true });

    window.location.href = "/";
  } catch (err) {
    console.error("error signing out:", err);
  }
}

export async function signIn(email, password) {
  const { isSignedIn, nextStep } = await amplifySignIn({
    username: email,
    password,
  });

  if (nextStep.signInStep === "CONFIRM_SIGN_UP") {
    window.location.href = "/confirm";
    return { isSignedIn: false, redirected: true };
  }

  if (isSignedIn) {
    const session = await fetchAuthSession();
    const accessToken = session.tokens.idToken.toString();
    localStorage.setItem("access_token", accessToken);
    window.location.href = "/";
  }

  return { isSignedIn, nextStep };
}
