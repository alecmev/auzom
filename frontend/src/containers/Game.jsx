import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import { userGameDataKeys } from '../utils';
import Component from '../utils/Component';

import NotFound from './NotFound';
import Loading from './Loading';

import Button from '../components/Button';
import Markdown from '../components/Markdown';

const gameSlugSelector = (state, props) => props.params.slug;
const gameSelector = createSelector(
  selectors.games,
  gameSlugSelector,
  (games, gameSlug) =>
    games.toList().find(x => x.get('slug') === gameSlug),
);

const gameIdSelector = createSelector(gameSelector, x => x && x.get('id'));
const gameTournamentsSelector = createSelector(
  selectors.tournaments,
  gameIdSelector,
  (tournaments, gameId) =>
    tournaments.toList().filter(x => x.get('gameId') === gameId),
);

const myOwnershipSelector = createSelector(
  selectors.userGames,
  selectors.myId,
  gameIdSelector,
  (x, myId, gameId) => x.toList().find(y =>
    y.get('userId') === myId &&
    y.get('gameId') === gameId &&
    y.get('nullifiedAt') === null
  ),
);

@connect(createStructuredSelector({
  isLoading: selectors.isLoading,
  myId: selectors.myId,
  amAdmin: selectors.amAdmin,
  tournaments: selectors.tournaments,

  gameSlug: gameSlugSelector,
  game: gameSelector,
  gameTournaments: gameTournamentsSelector,
  myOwnership: myOwnershipSelector,
}), actions)
export default class Game extends Component {
  static propTypes = {
    isLoading: PropTypes.bool.isRequired,
    myId: PropTypes.string,
    amAdmin: PropTypes.bool.isRequired,
    tournaments: ImmutablePropTypes.map,

    gameSlug: PropTypes.string.isRequired,
    game: ImmutablePropTypes.map,
    gameTournaments: ImmutablePropTypes.list,
    myOwnership: ImmutablePropTypes.map,

    loadGames: PropTypes.func.isRequired,
    loadTournaments: PropTypes.func.isRequired,
    createUserGame: PropTypes.func.isRequired,
    loadUserGame: PropTypes.func.isRequired,
    loadUserGames: PropTypes.func.isRequired,
    patchUserGame: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ myId, gameSlug, game }, prevProps) {
    if (
      !prevProps ||
      gameSlug !== prevProps.gameSlug ||
      (!game && prevProps.game)
    ) {
      this.props.loadGames({ slug: gameSlug });
    }

    if (!game) return;
    if (prevProps && game === prevProps.game) return;
    this.props.loadTournaments({ gameId: game.get('id') });
    if (!myId) return;
    this.props.loadUserGames({ userId: myId, gameId: game.get('id') });
  }

  handleAddClick = () => {
    if (!this.props.game) return;
    this.props.createUserGame({ gameId: this.props.game.get('id') });
  };
  handleStatusClick = () => {
    if (!this.props.myOwnership) return;
    this.props.loadUserGame(this.props.myOwnership.get('id'));
  };
  handleNullifyClick = () => {
    if (!this.props.myOwnership) return;
    if (!window.confirm('You sure?')) return;
    this.props.patchUserGame(this.props.myOwnership.get('id'), 'nullify');
  };

  render() {
    const {
      isLoading, myId, amAdmin, game: _game, gameTournaments, myOwnership,
    } = this.props;

    if (!_game) {
      return isLoading ?
        <div className={this.cni()}><Loading /></div> :
        <NotFound />;
    }

    const game = _game.toJS();

    return (
      <div className={this.cni()}>
        <div className={this.cn({ u: 'sectionDouble' })}>
          <div className={this.cn({
            u: 'sectionDoubleTwoThirdsLeft',
            d: 'left',
          })}>
            <div className={this.cn({ u: 'sectionMargined' })}>
              <img className={this.cn({ d: 'cover' })} src={game.cover} />
              <span className={this.cn({ d: 'info' })}>
                <div className={this.cn({ d: 'name' })}>{game.name}</div>
                {amAdmin &&
                  <div>
                    <Link
                      to={`/${game.slug}/_/settings`}
                      className={this.cn({ d: 'settings' })}
                    >
                      SETTINGS
                    </Link>
                  </div>
                }
                <Markdown
                  className={this.cn({ d: 'summary' })}
                  source={game.summary}
                />
                {!!gameTournaments.size && [
                  <div
                    className={this.cn({ d: 'tournamentsHeading' })}
                    key="tournamentsHeading"
                  >
                    Tournaments:
                  </div>,
                  <ul key="tournaments">{gameTournaments.toJS().map(x => (
                    <li key={x.id}>
                      <Link to={`/${game.slug}/${x.slug}`}>
                        {x.name}
                      </Link>
                    </li>
                  ))}</ul>,
                ]}
              </span>
            </div>
          </div>
          <div className={this.cn({
            u: 'sectionDoubleOneThirdRight',
            d: 'right',
          })}>
            <div className={this.cn({ u: 'sectionMargined' })}>
              <div className={this.cn({ d: 'rightInner' })}>
                {myId && !myOwnership && <div>
                  <Button
                    className={this.cn({ d: 'addButton' })}
                    text="add to my account" type="important" size="large"
                    onClick={this.handleAddClick}
                  />
                  <div className={this.cn({ d: 'disclaimer' })}>

                    After clicking the button above, all you&apos;ll have to do
                    is join our verification game server and write a short code
                    in the in-game chat, and that&apos;s it, we&apos;ll take
                    care of the rest.

                  </div>
                </div>}
                {myId && myOwnership && !myOwnership.get('verifiedAt') && <div>
                  <div className={this.cn({ d: 'warning' })}>
                    Verification in progress!
                  </div>
                  <div className={this.cn({ d: 'disclaimer' })}>
                    Join our
                    <a
                      className={this.cn({ d: 'link' })}
                      href="http://battlelog.battlefield.com/bf4/servers/show/pc/80b0676e-e188-4c75-8b10-89c2ea56155e/"
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      verification game server
                    </a>
                    and write
                    <code className={this.cn({ d: 'token' })}>
                      {myOwnership.get('token')}
                    </code>
                    in the in-game chat.
                  </div>
                  <Button
                    className={this.cn({ d: 'statusButton' })}
                    text="check status" type="important" size="large"
                    onClick={this.handleStatusClick}
                  />
                  <div className={this.cn({
                    d: 'nullifyButtonWrapper',
                    m: 'right',
                  })}>
                    <Button
                      text="cancel the process" type="important" size="small"
                      onClick={this.handleNullifyClick}
                    />
                  </div>
                </div>}
                {myId && myOwnership && myOwnership.get('verifiedAt') && <div>
                  <div className={this.cn({ d: 'warning' })}>
                    Noice, you have this game!
                  </div>
                  <div className={this.cn({ d: 'data' })}>
                    <div key="name">
                      <div className={this.cn({ d: 'dataKey' })}>
                        Name / Link
                      </div>
                      <div className={this.cn({ d: 'dataValue' })}>
                        <a href={myOwnership.get('link')}>
                          {myOwnership.get('name')}
                        </a>
                      </div>
                    </div>
                    {[...myOwnership.get('data').entries()].map(x =>
                      <div key={x[0]}>
                        <div className={this.cn({ d: 'dataKey' })}>
                          {userGameDataKeys[x[0]] || x[0]}
                        </div>
                        <div className={this.cn({ d: 'dataValue' })}>
                          {x[1]}
                        </div>
                      </div>
                    )}
                  </div>
                  <div className={this.cn({ d: 'nullifyButtonWrapper' })}>
                    <Button
                      text="remove from my account" type="important"
                      size="small" onClick={this.handleNullifyClick}
                    />
                  </div>
                </div>}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
