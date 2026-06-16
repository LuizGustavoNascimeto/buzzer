import "./DesktopSidebar.css";
import Search from "../components/Search";
import TrendingSection from "../components/TrendingsSection";
import SuggestedUsersSection from "../components/SuggestedUsersSection";
import JoinSection from "../components/JoinSection";

export default function DesktopSidebar(props) {
  const trendings = [
    { hashtag: "100DaysOfCloud", count: 2053 },
    { hashtag: "CloudProject", count: 8253 },
    { hashtag: "AWS", count: 9053 },
    { hashtag: "FreeWillyReboot", count: 7753 },
  ];

  const users = [{ display_name: "Andrew Brown", handle: "andrewbrown" }];

  return (
    <section>
      <Search />
      {props.handle ? (
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
