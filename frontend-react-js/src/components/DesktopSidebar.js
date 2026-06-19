import "./DesktopSidebar.css";
import Search from "../components/Search";
import TrendingSection from "../components/TrendingsSection";
import SuggestedUsersSection from "../components/SuggestedUsersSection";
import JoinSection from "../components/JoinSection";
import { useAuth } from "../hooks/auth/useAuth";

export default function DesktopSidebar(props) {
  const { data: user, isLoading: userLoading } = useAuth();
  const trendings = [
    { hashtag: "100DaysOfCloud", count: 2053 },
    { hashtag: "CloudProject", count: 8253 },
    { hashtag: "AWS", count: 9053 },
    { hashtag: "FreeWillyReboot", count: 7753 },
  ];

  const users = [{ display_name: "Andrew Brown", handle: "andrewbrown" }];

  const logged = () => {
    if (user) {
      if (user.handle) {
        return true;
      }
    }
    return false;
  };

  return (
    <section>
      <Search />
      {logged() ? (
        <>
          <TrendingSection trendings={trendings} />;
          <SuggestedUsersSection users={users} />
        </>
      ) : (
        <JoinSection />
      )}

      <footer>
        <a href="/about">About</a>
        <a href="/terms_of_service">Terms of Service</a>
        <a href="/privacy_policy">Privacy Policy</a>
      </footer>
    </section>
  );
}
