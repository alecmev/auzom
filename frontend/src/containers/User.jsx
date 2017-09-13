import moment from 'moment';
import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';
import { Link } from 'react-router';

import * as actions from '../actions';
import * as selectors from '../selectors';
import * as utils from '../utils';
import Component from '../utils/Component';

import NotFound from './NotFound';

import Button from '../components/Button';

const userIdSelector = (state, props) => props.params.id;
const userSelector = createSelector(
  selectors.users,
  userIdSelector,
  utils.get,
);

const mshipsSelector = createSelector(
  selectors.userTeams,
  userIdSelector,
  selectors.teams,
  (userTeams, userId, teams) =>
    userTeams.toList().filter(x =>
      x.get('userId') === userId &&
      x.get('leftAt') === null,
    ).sort((a, b) => {
      const leadership = +b.get('isLeader') - +a.get('isLeader');
      if (leadership !== 0) return leadership;
      const teamA = teams.get(a.get('teamId'));
      const teamB = teams.get(b.get('teamId'));
      if (!teamA || !teamB) return a.get('id') - b.get('id');
      return teamA.get('name').localeCompare(teamB.get('name'));
    }),
);

const pastMshipsSelector = createSelector(
  selectors.userTeams,
  userIdSelector,
  (userTeams, userId) =>
    userTeams.toList().filter(x =>
      x.get('userId') === userId &&
      x.get('leftAt') !== null,
    ).sort((a, b) => a.get('id') - b.get('id')),
);

const userGamesSelector = createSelector(
  selectors.userGames,
  userIdSelector,
  (x, userId) => x.toList().filter(y =>
    y.get('userId') === userId &&
    y.get('verifiedAt') !== null &&
    y.get('nullifiedAt') === null,
  ),
);

@connect(createStructuredSelector({
  myId: selectors.myId,
  amAdmin: selectors.amAdmin,
  teams: selectors.teams,
  games: selectors.games,

  userId: userIdSelector,
  user: userSelector,
  mships: mshipsSelector,
  pastMships: pastMshipsSelector,
  userGames: userGamesSelector,
}), actions)
export default class User extends Component {
  static propTypes = {
    myId: PropTypes.string,
    amAdmin: PropTypes.bool,
    teams: ImmutablePropTypes.map,
    games: ImmutablePropTypes.map,

    userId: PropTypes.string,
    user: ImmutablePropTypes.map,
    mships: ImmutablePropTypes.list,
    pastMships: ImmutablePropTypes.list,
    userGames: ImmutablePropTypes.list,

    loadUser: PropTypes.func.isRequired,
    updateUser: PropTypes.func.isRequired,
    loadUserTeams: PropTypes.func.isRequired,
    patchUserGame: PropTypes.func.isRequired,
    loadUserGame: PropTypes.func.isRequired,
    loadUserGames: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ userId, user, mships, pastMships }, prevProps) {
    if (
      !prevProps ||
      userId !== prevProps.userId ||
      (!user && prevProps.user)
    ) {
      this.props.loadUser(userId);
      this.props.loadUserTeams(
        { userId }, // , leftAt: '%00'
        mships.concat(pastMships).reduce((r, x) => r.concat(x.get('id')), []),
        ['teamId'],
      );
      this.props.loadUserGames({ userId }, undefined, ['gameId']);
    }
  }

  handleGiveAdmin = () => {
    this.props.updateUser(this.props.userId, { isAdmin: true });
  };

  handleRemoveAdmin = () => {
    this.props.updateUser(this.props.userId, { isAdmin: false });
  };

  handleUpdateRequestClick = (_, id) => this.props.patchUserGame(id, 'update');
  handleRefreshClick = (_, id) => this.props.loadUserGame(id);

  render() {
    const {
      myId, amAdmin, teams, games, user: _user, mships, pastMships, userGames,
    } = this.props;
    if (!_user) {
      return null;
    }

    if (_user === utils.NOT_FOUND) {
      return <NotFound />;
    }

    const user = _user.toJS();
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined' })}>
          User: {user.nickname}{amAdmin &&
            <span className={this.cn({ d: 'secondary' })}>
              &nbsp;/ {user.email}
            </span>
          }
          <div>
            Is admin: {user.isAdmin ?
              <span className={this.cn({ d: 'isAdmin' })}>yes</span> :
              <span className={this.cn({ d: 'notAdmin' })}>no</span>
            }
            {amAdmin && !user.isAdmin &&
              <Button
                text="give admin"
                type="important"
                size="small"
                onClick={this.handleGiveAdmin}
              />
            }
            {amAdmin && user.isAdmin &&
              <Button
                text="remove admin"
                type="important"
                size="small"
                onClick={this.handleRemoveAdmin}
              />
            }
          </div>
          {!!mships.size && [
            <div key="mships-heading">
              Memberships
              <span className={this.cn({ d: 'secondary' })}>
                &nbsp;[{mships.size}]
              </span>
            </div>,
            <ul key="mships">{mships.toJS().map((mship) => {
              const team = teams.get(mship.teamId);
              return (
                <li key={mship.id}>
                  {team && <Link to={`/teams/${team.get('id')}`}>
                    {`[${team.get('abbr')}] ${team.get('name')}`}
                  </Link>}
                  {mship.isLeader &&
                    <span className={this.cn({ d: 'secondary' })}>
                      &nbsp;/ leader of
                    </span>
                  }
                </li>
              );
            })}</ul>,
          ]}
          {!!pastMships.size && [
            <div key="pastMships-heading">
              Past memberships
              <span className={this.cn({ d: 'secondary' })}>
                &nbsp;[{pastMships.size}]
              </span>
            </div>,
            <ul key="pastMships">{pastMships.toJS().map((mship) => {
              const team = teams.get(mship.teamId);
              return (
                <li key={mship.id}>
                  {team && <Link to={`/teams/${team.get('id')}`}>
                    {`[${team.get('abbr')}] ${team.get('name')}`}
                  </Link>}
                  <span className={this.cn({ d: 'secondary' })}>{` /
                    ${moment(mship.createdAt).format('YYYY-MM-DD')} ~
                    ${moment(mship.leftAt).format('YYYY-MM-DD')} /
                    ${mship.kickedBy ? 'kicked' : 'left'}
                  `}</span>
                </li>
              );
            })}</ul>,
          ]}
          {!!userGames.size && [
            <div key="userGames-heading">User&apos;s games</div>,
            <ul key="userGames">{userGames.toJS().map((x) => {
              const game = games.get(x.gameId);
              return (
                <li key={x.id}>
                  <Link to={`/${game && game.get('slug')}`}>
                    {game ? `${game.get('name')}` : `Game ${x.gameId}`}
                  </Link>
                  : <a href={x.link}>{x.name}</a>
                  <span className={this.cn({ d: 'secondary' })}>{` /
                    updated ${moment(x.dataUpdatedAt).fromNow()}
                  `}</span>
                  {x.dataUpdateRequestedAt ? <span>
                    <span
                      className={this.cn({ d: 'secondary' })}
                    >{` / update requested
                      ${moment(x.dataUpdateRequestedAt).fromNow()}
                    `}</span>
                    <Button
                      text="refresh"
                      type="important"
                      size="small"
                      meta={x.id}
                      onClick={this.handleRefreshClick}
                    />
                  </span> : <Button
                    text="request update"
                    type="important"
                    size="small"
                    meta={x.id}
                    onClick={this.handleUpdateRequestClick}
                  />}
                </li>
              );
            })}</ul>,
          ]}
          {!userGames.size && myId === user.id && <Button
            text="add battlefield 4..."
            type="important"
            size="large"
            href="/battlefield-4"
          />}
        </div>
      </div>
    );
  }
}
