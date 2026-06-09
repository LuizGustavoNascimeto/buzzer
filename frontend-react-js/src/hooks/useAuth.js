import { useCallback, useEffect, useState } from "react";
import {
  getCurrentUser,
  fetchUserAttributes,
  fetchAuthSession,
  signOut as amplifySignOut,
} from "aws-amplify/auth";
import { Hub } from "aws-amplify/utils";

export function useAuth() {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Única fonte da verdade — busca usuário E token juntos
  const checkAuth = useCallback(async (forceRefresh = false) => {
    try {
      setLoading(true);
      setError(null);

      // 1. Verifica se há usuário logado
      await getCurrentUser();

      // 2. Busca atributos e token em paralelo
      const [attributes, session] = await Promise.all([
        fetchUserAttributes(),
        fetchAuthSession({ forceRefresh }),
      ]);

      const idToken = session.tokens?.idToken?.toString();

      if (!idToken) throw new Error("Token não encontrado");

      setUser({
        display_name: attributes.preferred_username,
        handle: attributes.name,
      });
      setToken(idToken);

      return idToken;
    } catch (err) {
      console.error(err);
      setError(err);
      setUser(null);
      setToken(null);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const signOut = useCallback(async () => {
    try {
      await amplifySignOut({ global: true });
      setUser(null); // limpa estado em memória
      setToken(null);
      window.location.href = "/";
    } catch (err) {
      console.error("error signing out:", err);
    }
  }, []);

  // Checagem inicial
  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  // Escuta eventos do Amplify
  useEffect(() => {
    const unsubscribe = Hub.listen("auth", ({ payload }) => {
      if (payload.event === "tokenRefresh") {
        checkAuth(); // re-sincroniza tudo
      }
      if (payload.event === "tokenRefresh_failure") {
        setToken(null);
        setUser(null);
        window.location.href = "/login";
      }
    });

    return unsubscribe;
  }, [checkAuth]);

  return {
    checkAuth,
    refreshToken: () => checkAuth(true), // forceRefresh explícito
    token,
    user,
    loading,
    error,
    refreshUser: checkAuth,
    signOut,
  };
}
