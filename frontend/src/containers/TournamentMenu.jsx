import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Link, IndexLink } from 'react-router';

import * as utils from '../utils';
import Component from '../utils/Component';

export default class TournamentMenu extends Component {
  static propTypes = {
    isLoading: PropTypes.bool.isRequired,
    gameSlug: PropTypes.string.isRequired,
    game: ImmutablePropTypes.map,
    tournamentPath: PropTypes.string.isRequired,
    tournament: ImmutablePropTypes.map,
    seasonPath: PropTypes.string.isRequired,
    season: ImmutablePropTypes.map,
    myParticipation: ImmutablePropTypes.map,
    myApplication: ImmutablePropTypes.map,
    teams: ImmutablePropTypes.map,
  };

  render() {
    const {
      isLoading, gameSlug, game, tournament, season: _season,
      tournamentPath, seasonPath, myParticipation, myApplication, teams,
    } = this.props;

    if ((!game || !tournament) && !isLoading) {
      return null;
    }

    let season;
    let seasonName = 'LOADING';
    let status = 'LOADING';
    let statusColor = 'gray';
    let statusBold = false;
    let seasonPathIfKnown = tournamentPath;
    if (_season) {
      season = _season.toJS();
      seasonName = season.name;
      seasonPathIfKnown = seasonPath;
      const now = new Date();
      if (season.endedAt && new Date(season.endedAt) <= now) {
        status = 'ENDED';
        statusColor = 'red';
      } else if (
        season.signupsClosedAt &&
        new Date(season.signupsClosedAt) <= now
      ) {
        status = 'IN PROGRESS';
        statusColor = 'yellow';
      } else if (
        season.signupsOpenedAt &&
        new Date(season.signupsOpenedAt) <= now
      ) {
        status = 'SIGNUPS OPEN';
        statusColor = 'green';
        statusBold = true;
      } else {
        status = 'ANNOUNCED';
      }
    } else if (!isLoading) {
      seasonName = 'NO ACTIVE SEASON / NO SEASON SELECTED';
      status = '';
    }

    let myTeam;
    let warnTeam;
    if (
      myParticipation &&
      myApplication &&
      new Date(myParticipation.get('createdAt')) >
      new Date(myApplication.get('createdAt'))
    ) {
      myTeam = teams.get(myParticipation.get('teamId'));
      warnTeam = myParticipation.get('leftAt') && (
        !myParticipation.get('isDone') ||
        myParticipation.get('kickedBy')
      );
    } else if (myApplication) {
      myTeam = teams.get(myApplication.get('teamId'));
      warnTeam = myApplication.get('decision') === false;
    }

    return (
      <div className={this.cni()}>
        <div className={this.cn({ d: 'title' })}>
          <div className={this.cn({ d: 'logo' })}>
            <div
              className={this.cn({ d: 'logoBackground' })}
              style={{
                backgroundImage: `url(${utils.https(tournament.get('logo'))})`,
              }}
            />
            <div className={this.cn({ d: 'logoTag' })}>
              {!tournament.get('logoHasText') && tournament.get('abbr')}
            </div>
          </div>
          <div className={this.cn({ d: 'tournament' })}>
            <Link
              to={`/${gameSlug}`}
              className={this.cn({ d: 'tournament-game' })}
            >
              {game ? game.get('name') : 'LOADING'}
            </Link>
            <div className={this.cn({ d: 'tournament-name' })}>
              {tournament ? tournament.get('name') : 'LOADING'}
            </div>
            <div className={this.cn({ d: 'tournament-season' })}>
              {seasonName}
              <span className={this.cn({
                d: 'tournament-status',
                m: [statusColor, statusBold && 'bold'],
              })}>
                {status}
              </span>
            </div>
          </div>
        </div>
        <div className={this.cn({ d: 'entries' })}>
          <IndexLink
            to={`${seasonPathIfKnown}`}
            activeClassName="is-active"
            className={this.cn({ d: 'entry' })}
          >
            home
          </IndexLink>
          <Link
            to={`${seasonPathIfKnown}/rules`}
            activeClassName="is-active"
            className={this.cn({ d: 'entry' })}
          >
            rules
          </Link>
          <Link
            to={`${seasonPathIfKnown}/participants`}
            activeClassName="is-active"
            className={this.cn({ d: 'entry' })}
          >
            participants
          </Link>
          <Link
            to={`${seasonPathIfKnown}/brackets`}
            activeClassName="is-active"
            className={this.cn({ d: 'entry' })}
          >
            brackets
          </Link>
          {season && season.youtubePlaylist && <Link
            to={`${seasonPathIfKnown}/videos`}
            activeClassName="is-active"
            className={this.cn({ d: 'entry' })}
          >
            videos
          </Link>}
          <Link
            to={`${tournamentPath}/seasons`}
            activeClassName="is-active"
            className={this.cn({ d: 'entry' })}
          >
            seasons
          </Link>
          <div className={this.cn({ d: 'entriesFill' })} />
          {myTeam && <Link
            to={`${seasonPathIfKnown}/participants/${myTeam.get('id')}`}
            activeClassName="is-active"
            className={this.cn({
              d: 'entry',
              m: [warnTeam ? 'warn' : 'highlight'],
            })}
          >
            {warnTeam && '❗'} my team {warnTeam && '❗'}
          </Link>}
        </div>
      </div>
    );
  }
}
