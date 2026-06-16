import "./App.css";

import {QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./lib/queryClient";
import process from "process";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import ConfirmationPage from "./pages/ConfirmationPage";
import HomeFeedPage from "./pages/HomeFeedPage";
import MessageGroupPage from "./pages/MessageGroupPage";
import MessageGroupsPage from "./pages/MessageGroupsPage";
import RecoverPage from "./pages/RecoverPage";
import SigninPage from "./pages/SigninPage";
import SignupPage from "./pages/SignupPage";
import UserFeedPage from "./pages/UserFeedPage";

import { Amplify } from "aws-amplify";
import NotificationsFeedPage from "./pages/NotificationsFeedPage";

Amplify.configure({
  Auth: {
    Cognito: {
      userPoolId: process.env.REACT_APP_AWS_USER_POOLS_ID,
      userPoolClientId: process.env.REACT_APP_CLIENT_ID,
    },
  },
});

const router = createBrowserRouter([
  {
    path: "/",
    element: <HomeFeedPage />,
  },
  {
    path: "/notifications",
    element: <NotificationsFeedPage />,
  },
  {
    path: "/@:handle",
    element: <UserFeedPage />,
  },
  {
    path: "/messages",
    element: <MessageGroupsPage />,
  },
  {
    path: "/messages/:group_id",
    element: <MessageGroupPage />,
  },
  {
    path: "/signup",
    element: <SignupPage />,
  },
  {
    path: "/signin",
    element: <SigninPage />,
  },
  {
    path: "/confirm",
    element: <ConfirmationPage />,
  },
  {
    path: "/forgot",
    element: <RecoverPage />,
  },
]);

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  );
}

export default App;
