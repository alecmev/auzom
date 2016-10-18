import { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import {
  applyRouterMiddleware, browserHistory,
  Router as ReactRouter, Route, IndexRoute, Redirect,
} from 'react-router';
import { useScroll } from 'react-router-scroll';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import store from '../store';

import About from './About';
import AttentionRequestNew from './AttentionRequestNew';
import AttentionRequests from './AttentionRequests';
import Bracket from './Bracket';
import BracketNew from './BracketNew';
import BracketMapNew from './BracketMapNew';
import BracketNewSwissRound from './BracketNewSwissRound';
import BracketPrepare from './BracketPrepare';
import BracketSettings from './BracketSettings';
import Game from './Game';
import GameMapNew from './GameMapNew';
import GameNew from './GameNew';
import Games from './Games';
import GameSettings from './GameSettings';
import Home from './Home';
import Loading from './Loading';
import Login from './Login';
import Match from './Match';
import MatchComments from './MatchComments';
import MatchDetails from './MatchDetails';
import MatchMaps from './MatchMaps';
import MatchPrepare from './MatchPrepare';
import MatchReportNew from './MatchReportNew';
import MatchReports from './MatchReports';
import MatchSettings from './MatchSettings';
import NotFound from './NotFound';
import Page from './Page';
import PasswordReset from './PasswordReset';
import SeasonApply from './SeasonApply';
import SeasonBrackets from './SeasonBrackets';
import SeasonHome from './SeasonHome';
import SeasonNew from './SeasonNew';
import SeasonNewsItem from './SeasonNewsItem';
import SeasonNewsItemEditor from './SeasonNewsItemEditor';
import SeasonParticipant from './SeasonParticipant';
import SeasonParticipants from './SeasonParticipants';
import SeasonRules from './SeasonRules';
import SeasonRulesEditor from './SeasonRulesEditor';
import SeasonSettings from './SeasonSettings';
import SeasonVideos from './SeasonVideos';
import Signup from './Signup';
import StageNew from './StageNew';
import StageSettings from './StageSettings';
import Team from './Team';
import TeamNew from './TeamNew';
import Teams from './Teams';
import TeamSettings from './TeamSettings';
import Tournament from './Tournament';
import TournamentNew from './TournamentNew';
import TournamentSeasons from './TournamentSeasons';
import TournamentSettings from './TournamentSettings';
import User from './User';
import UserSettings from './UserSettings';

if (!__DEV__) {
  window.ga = window.ga || function ga(...args) {
    window.ga.q = window.ga.q || [];
    window.ga.q.push(args);
  };
  window.ga.l = +new Date();
  window.ga('create', 'UA-78826698-1', 'auto');

  let previousPath = null;
  browserHistory.listen((loc) => {
    const currentPath = loc.path || (loc.pathname + loc.search);
    if (previousPath === currentPath) return;
    window.ga('set', 'page', currentPath);
    window.ga('send', 'pageview');
    previousPath = currentPath;
  });

  let previousUserId = null;
  store.subscribe(() => {
    const currentUserId = selectors.myId(store.getState());
    if (previousUserId === currentUserId) return;
    window.ga('userId', currentUserId);
    previousUserId = currentUserId;
  });
}

@connect(createStructuredSelector({
  isRehydrated: selectors.isRehydrated,
  myId: selectors.myId,
  amAdmin: selectors.amAdmin,
}), actions)
export default class Router extends Component {
  static propTypes = {
    isRehydrated: PropTypes.bool.isRequired,
    myId: PropTypes.string,
    amAdmin: PropTypes.bool.isRequired,
    createSession: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  shouldComponentUpdate(nextProps) {
    return nextProps.isRehydrated !== this.props.isRehydrated;
  }

  verify = (nextState) => {
    this.props.createSession({ token: nextState.params.token }, () => {
      browserHistory.push('/');
      this.props.messagePush('email verified');
    });
  };

  passwordReset = (nextState) => {
    this.props.createSession({ token: nextState.params.token }, () => {
      browserHistory.push('/settings');
      this.props.messagePush('logged in, now change your password');
    });
  };

  redirectIfLoggedIn = (nextState, replace) => {
    if (this.props.myId) {
      replace('/');
    }
  };

  redirectIfLoggedOut = (nextState, replace) => {
    if (!this.props.myId) {
      replace('/login');
    }
  };

  redirectIfNotAdmin = (nextState, replace) => {
    // TODO: or maybe just 404
    if (!this.props.amAdmin) {
      this.props.messagePush(
        `you must be an admin to access ${nextState.location.pathname}`, true,
      );
      replace('/');
    }
  };

  render() {
    if (!this.props.isRehydrated) {
      return <Loading />;
    }

    let kitchenSink;
    if (__DEV__) {
      const KitchenSink = require('./KitchenSink').default;

      kitchenSink = (
        <Route
          path="kitchen-sink"
          component={KitchenSink}
        />
      );
    }

    return (
      <ReactRouter
        history={browserHistory}
        render={applyRouterMiddleware(useScroll(
          (_, { location }) => !location.hash
        ))}
      >
        <Route path="/" component={Page}>
          <IndexRoute component={Home} />
          {kitchenSink}
          {__DEV__ && <Route
            path="loading"
            component={Loading}
          />}
          <Route
            path="signup"
            component={Signup}
            onEnter={this.redirectIfLoggedIn}
          />
          <Route
            path="verify/:token"
            /* TODO: this needs a placeholder component */
            onEnter={this.verify}
          />
          <Route
            path="login"
            component={Login}
            onEnter={this.redirectIfLoggedIn}
          />
          <Route path="password-reset">
            <IndexRoute
              component={PasswordReset}
              onEnter={this.redirectIfLoggedIn}
            />
            <Route
              path=":token"
              onEnter={this.passwordReset}
            />
          </Route>
          <Route
            path="settings"
            component={UserSettings}
            onEnter={this.redirectIfLoggedOut}
          />
          <Route
            path="about"
            component={About}
          />
          <Route path="users">
            <Route
              path=":id"
              component={User}
            />
          </Route>
          <Route path="teams">
            <IndexRoute component={Teams} />
            <Route
              path="new"
              component={TeamNew}
              onEnter={this.redirectIfLoggedOut}
            />
            <Route
              path=":id"
              component={Team}
            />
            <Route
              path=":id/settings"
              component={TeamSettings}
            />
          </Route>
          <Route path="games">
            <IndexRoute component={Games} />
            <Route
              path="new"
              component={GameNew}
              onEnter={this.redirectIfNotAdmin}
            />
          </Route>
          <Route path="game-maps">
            <Route
              path="new"
              component={GameMapNew}
              onEnter={this.redirectIfNotAdmin}
            />
          </Route>
          <Route path="tournaments">
            <Route
              path="new"
              component={TournamentNew}
              onEnter={this.redirectIfNotAdmin}
            />
          </Route>
          <Route path="seasons">
            <Route
              path="new"
              component={SeasonNew}
              onEnter={this.redirectIfNotAdmin}
            />
          </Route>
          <Route
            path=":slug"
            component={Game}
          />
          <Route
            path=":slug/_/settings"
            component={GameSettings}
            onEnter={this.redirectIfNotAdmin}
          />
          <Route
            path=":gameSlug/:tournamentSlug"
            component={Tournament}
          >
            <Route
              path="seasons"
              component={TournamentSeasons}
            />
            <Route
              path="_/settings"
              component={TournamentSettings}
              onEnter={this.redirectIfNotAdmin}
            />
            <Route path=":seasonSlug">
              <IndexRoute component={SeasonHome} />
              <Route
                path="_/settings"
                component={SeasonSettings}
                onEnter={this.redirectIfNotAdmin}
              />
              <Route
                path="apply"
                component={SeasonApply}
              />
              <Route path="rules">
                <IndexRoute component={SeasonRules} />
                <Route
                  path="edit"
                  component={SeasonRulesEditor}
                  onEnter={this.redirectIfNotAdmin}
                />
              </Route>
              <Route path="participants">
                <IndexRoute component={SeasonParticipants} />
                <Route
                  path=":teamId"
                  component={SeasonParticipant}
                />
              </Route>
              <Route path="news">
                <Route
                  path="new"
                  component={
                    props => <SeasonNewsItemEditor
                      {...props}
                      isNew
                    />
                  }
                  onEnter={this.redirectIfNotAdmin}
                />
                <Route path=":newsItemId">
                  <IndexRoute component={SeasonNewsItem} />
                  <Route
                    path="edit"
                    component={SeasonNewsItemEditor}
                    onEnter={this.redirectIfNotAdmin}
                  />
                </Route>
              </Route>
              <Route
                path="videos"
                component={SeasonVideos}
              />
              <Redirect from="overview/*" to="brackets/*" />
              <Route
                path="brackets"
                component={SeasonBrackets}
              >
                <Route
                  path="_/new-stage"
                  component={StageNew}
                  onEnter={this.redirectIfNotAdmin}
                />
                <Route path=":stageSlug">
                  <Route
                    path="_"
                    onEnter={this.redirectIfNotAdmin}
                  >
                    <Route
                      path="settings"
                      component={StageSettings}
                    />
                    <Route
                      path="new-bracket"
                      component={BracketNew}
                    />
                  </Route>
                  <Route path=":bracketSlug">
                    <IndexRoute component={Bracket} />
                    <Route
                      path="_"
                      onEnter={this.redirectIfNotAdmin}
                    >
                      <Route
                        path="new-map"
                        component={BracketMapNew}
                      />
                      <Route
                        path="new-swiss-round"
                        component={BracketNewSwissRound}
                      />
                      <Route
                        path="prepare"
                        component={BracketPrepare}
                      />
                      <Route
                        path="settings"
                        component={BracketSettings}
                      />
                    </Route>
                    <Route
                      path="matches/:matchId"
                      component={Match}
                    >
                      <IndexRoute component={MatchComments} />
                      <Route
                        path="maps"
                        component={MatchMaps}
                      />
                      <Route
                        path="details"
                        component={MatchDetails}
                      />
                      <Route path="reports">
                        <IndexRoute component={MatchReports} />
                        <Route
                          path="new"
                          component={MatchReportNew}
                        />
                      </Route>
                      <Route path="attention-requests">
                        <IndexRoute component={AttentionRequests} />
                        <Route
                          path="new"
                          component={AttentionRequestNew}
                        />
                      </Route>
                      <Route
                        path="_"
                        onEnter={this.redirectIfNotAdmin}
                      >
                        <Route
                          path="settings"
                          component={MatchSettings}
                        />
                        <Route
                          path="prepare"
                          component={MatchPrepare}
                        />
                      </Route>
                    </Route>
                  </Route>
                </Route>
              </Route>
            </Route>
          </Route>
          <Route path="*" component={NotFound} />
        </Route>
      </ReactRouter>
    );
  }
}
