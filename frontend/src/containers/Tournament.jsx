import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import * as utils from '../utils';
import Component from '../utils/Component';

import Loading from './Loading';
import NotFound from './NotFound';
import TournamentMenu from './TournamentMenu';

const gameSlugSelector = (state, props) => props.params.gameSlug;
const gameSelector = createSelector(
  selectors.games,
  gameSlugSelector,
  (games, slug) => games.toList().find(x => x.get('slug') === slug),
);

const tournamentSlugSelector = (state, props) => props.params.tournamentSlug;
const tournamentPathSelector = (_, { params }) =>
  `/${params.gameSlug}/${params.tournamentSlug}`;
const tournamentSelector = createSelector(
  selectors.tournaments,
  createSelector(gameSelector, x => x && x.get('id')),
  tournamentSlugSelector,
  (tournaments, gameId, slug) => tournaments.toList().find(x =>
    x.get('gameId') === gameId &&
    x.get('slug') === slug
  ),
);

const tournamentSeasonsSelector = createSelector(
  selectors.seasons,
  createSelector(tournamentSelector, x => x && x.get('id')),
  (x, tournamentId) =>
    tournamentId && x.toList()
      .filter(y => y.get('tournamentId') === tournamentId)
      .sortBy(y => y.get('publishedAt')),
);

const seasonSlugSelector = (_, props) => props.params.seasonSlug;
const seasonPathSelector = (_, { params }) =>
  `/${params.gameSlug}/${params.tournamentSlug}/${params.seasonSlug}`;
const seasonSelector = createSelector(
  tournamentSeasonsSelector,
  seasonSlugSelector,
  (x, slug) => x && slug && x.find(y => y.get('slug') === slug),
);

const membershipsSelector = createSelector(
  selectors.userTeams,
  selectors.myId,
  (userTeams, myId) =>
    myId && userTeams.toList().filter(x =>
      x.get('userId') === myId &&
      x.get('leftAt') === null
    ),
);

const leadershipsSelector = createSelector(
  membershipsSelector,
  x => x && x.filter(y => y.get('isLeader')),
);

const membershipTeamIdsSelector = createSelector(
  membershipsSelector,
  x => x && x.reduce((r, y) => r.concat(y.get('teamId')), []),
);

const participantsSelector = createSelector(
  selectors.teamSeasons,
  seasonSelector,
  selectors.teams,
  (teamSeasons, season, teams) =>
    season && teamSeasons.toList().filter(x =>
      x.get('seasonId') === season.get('id')
    ).sort((a, b) => {
      const ta = teams.get(a.get('teamId'));
      if (ta) {
        const tb = teams.get(b.get('teamId'));
        if (tb) return ta.get('name').localeCompare(tb.get('name'));
      }

      return a.get('teamId') - b.get('teamId');
    }),
);

const myParticipationSelector = createSelector(
  participantsSelector,
  membershipTeamIdsSelector,
  (x, ids) => x && ids && x.filter(
    y => ids.includes(y.get('teamId'))
  ).sort(
    (a, b) => b.get('createdAt').localeCompare(a.get('createdAt'))
  ).first(),
);

const applicationsSelector = createSelector(
  selectors.teamSeasonRequests,
  seasonSelector,
  selectors.teams,
  (teamSeasonRequests, season, teams) =>
    season && teamSeasonRequests.toList().filter(x =>
      x.get('seasonId') === season.get('id')
    ).sort((a, b) => {
      const ta = teams.get(a.get('teamId'));
      if (ta) {
        const tb = teams.get(b.get('teamId'));
        if (tb) return ta.get('name').localeCompare(tb.get('name'));
      }

      return a.get('teamId') - b.get('teamId');
    }),
);

const myApplicationSelector = createSelector(
  applicationsSelector,
  membershipTeamIdsSelector,
  (x, ids) => x && ids && x.filter(
    y => ids.includes(y.get('teamId'))
  ).sort(
    (a, b) => b.get('createdAt').localeCompare(a.get('createdAt'))
  ).first(),
);

const shortcuts = ['rules', 'participants', 'brackets', 'videos'];

@connect(createStructuredSelector({
  isLoading: selectors.isLoading,
  myId: selectors.myId,
  gameSlug: gameSlugSelector,
  game: gameSelector,
  tournamentSlug: tournamentSlugSelector,
  tournamentPath: tournamentPathSelector,
  tournament: tournamentSelector,
  tournamentSeasons: tournamentSeasonsSelector,
  seasonSlug: seasonSlugSelector,
  seasonPath: seasonPathSelector,
  season: seasonSelector,
  memberships: membershipsSelector,
  leaderships: leadershipsSelector,
  participants: participantsSelector,
  myParticipation: myParticipationSelector,
  applications: applicationsSelector,
  myApplication: myApplicationSelector,
  teams: selectors.teams, // for TournamentMenu
}), actions)
export default class Tournament extends Component {
  static propTypes = {
    isLoading: PropTypes.bool.isRequired,
    myId: PropTypes.string,
    children: PropTypes.element,
    gameSlug: PropTypes.string,
    game: ImmutablePropTypes.map,
    tournamentSlug: PropTypes.string,
    tournamentPath: PropTypes.string.isRequired,
    tournament: ImmutablePropTypes.map,
    tournamentSeasons: ImmutablePropTypes.list,
    seasonSlug: PropTypes.string,
    seasonPath: PropTypes.string.isRequired,
    season: ImmutablePropTypes.map,
    memberships: ImmutablePropTypes.list,
    participants: ImmutablePropTypes.list,
    applications: ImmutablePropTypes.list,

    loadGames: PropTypes.func.isRequired,
    loadTournaments: PropTypes.func.isRequired,
    loadSeasons: PropTypes.func.isRequired,
    loadUserTeams: PropTypes.func.isRequired,
    loadTeamSeasons: PropTypes.func.isRequired,
    loadTeamSeasonRequests: PropTypes.func.isRequired,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({
    myId, gameSlug, game, tournamentSlug, tournamentPath, tournament,
    tournamentSeasons, seasonSlug, season, children, location,
    memberships, participants, applications,
  }, prevProps) {
    if (myId && (!prevProps || myId !== prevProps.myId)) {
      this.props.loadUserTeams(
        { userId: myId, leftAt: '%00' },
        memberships.reduce((r, x) => r.concat(x.get('id')), []),
        ['teamId'],
      );
    }

    if (
      !prevProps ||
      gameSlug !== prevProps.gameSlug ||
      (!game && prevProps.game)
    ) {
      this.props.loadGames({ slug: gameSlug });
    }

    if (!game) return;
    if (
      !prevProps ||
      game !== prevProps.game ||
      tournamentSlug !== prevProps.tournamentSlug ||
      (!tournament && prevProps.tournament)
    ) {
      this.props.loadTournaments(
        { gameId: game.get('id'), slug: tournamentSlug },
      );
    }

    if (!tournament) return;
    if (
      !prevProps ||
      tournament !== prevProps.tournament ||
      (!tournamentSeasons && prevProps.tournamentSeasons)
    ) {
      this.props.loadSeasons({ tournamentId: tournament.get('id') });
    }

    if (!season) {
      if (
        !tournamentSeasons || !tournamentSeasons.size || (
          location.pathname !== tournamentPath &&
          !shortcuts.includes(seasonSlug)
      )) return;

      const now = new Date();
      let currentSeason;
      tournamentSeasons.forEach((x) => {
        if (new Date(x.get('publishedAt')) > now) {
          if (!currentSeason) currentSeason = x;
          return false;
        }

        currentSeason = x;
        return true;
      });

      const url = `${tournamentPath}/${currentSeason.get('slug')}`;
      if (!shortcuts.includes(seasonSlug)) {
        this.context.router.replace(url);
      } else {
        this.context.router.replace(
          `${url}${location.pathname.replace(tournamentPath, '')}`,
        );
      }

      return;
    }
    if (!prevProps || season !== prevProps.season) {
      this.props.loadTeamSeasons(
        { seasonId: season.get('id') }, // , leftAt: '%00'
        participants.reduce((r, x) => r.concat(x.get('id')), []),
        ['teamId'],
      );
      this.props.loadTeamSeasonRequests(
        { seasonId: season.get('id') }, // , decision: '%00'
        applications.reduce((r, x) => r.concat(x.get('id')), []),
        ['teamId'],
      );
    }
  }

  render() {
    const { isLoading, game, tournament, season } = this.props;

    if (!game || !tournament) {
      return isLoading ?
        <div className={this.cni()}><Loading /></div> :
          <NotFound />;
    }

    const childProps = utils.cloneProps(this.props);
    const children = React.Children.map(
      this.props.children, x => React.cloneElement(x, childProps),
    );

    const sponsors = [];
    if (season) {
      for (const sponsor of season.get('sponsors').split('\n')) {
        const sponsorItems = sponsor.trim().split(';');
        if (sponsorItems.length < 2) continue;
        sponsors.push({
          title: sponsorItems[0],
          logo: sponsorItems[1],
          scale: sponsorItems[2] || 50,
          url: sponsorItems[3] || '',
        });
      }
    }

    return (
      <div className={this.cni()}>
        <div
          className={this.cn({ u: 'sectionSingle', d: 'menu' })}
          style={{
            backgroundImage: `url(${utils.https(tournament.get('blur'))})`,
          }}
        >
          <div className={this.cn({ u: 'sectionMargined' })}>
            <TournamentMenu {...childProps} />
          </div>
        </div>
        {!!sponsors.length && <div
          className={this.cn({ u: 'sectionSingle', d: 'sponsors' })}
        >
          <div className={this.cn({
            u: 'sectionMargined',
            d: 'sponsorsInner',
          })}>
            {sponsors.map((x, i) => (
              <div
                key={i}
                className={this.cn({ d: 'sponsorsItemOuter' })}
              >
                <a
                  className={this.cn({ d: 'sponsorsItem' })}
                  style={{
                    backgroundImage: `url(${utils.https(x.logo)})`,
                    backgroundSize: `auto ${x.scale}%`,
                  }}
                  href={x.url}
                  title={x.title}
                />
              </div>
            ))}
          </div>
        </div>}
        {children}
      </div>
    );
  }
}
