import moment from 'moment';
import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import * as utils from '../utils';
import Component from '../utils/Component';

import Loading from './Loading';
import NotFound from './NotFound';

import Button from '../components/Button';
import Markdown from '../components/Markdown';

const newsItemsSelector = createSelector(
  selectors.newsItems,
  (_, props) => props.season,
  (x, season) =>
    season && x.toList().filter(y =>
      y.get('target') === 'season' &&
      y.get('targetId') === season.get('id')
    ).sort((a, b) => {
      const ap = a.get('publishedAt');
      const bp = b.get('publishedAt');
      if (ap === null && bp === null) return 0;
      else if (ap === null) return 1;
      else if (bp === null) return -1;
      return bp.localeCompare(ap);
    }),
);

@connect(createStructuredSelector({
  isLoading: selectors.isLoading,
  myId: selectors.myId,
  amAdmin: selectors.amAdmin,
  teams: selectors.teams,

  newsItems: newsItemsSelector,
}), actions)
export default class SeasonHome extends Component {
  static propTypes = {
    game: ImmutablePropTypes.map,
    tournamentPath: PropTypes.string,
    tournament: ImmutablePropTypes.map,
    seasonPath: PropTypes.string,
    season: ImmutablePropTypes.map,

    isLoading: PropTypes.bool.isRequired,
    myId: PropTypes.string,
    amAdmin: PropTypes.bool,
    teams: ImmutablePropTypes.map,

    newsItems: ImmutablePropTypes.list,

    memberships: ImmutablePropTypes.list,
    leaderships: ImmutablePropTypes.list,
    participants: ImmutablePropTypes.list,
    myParticipation: ImmutablePropTypes.map,
    applications: ImmutablePropTypes.list,
    myApplication: ImmutablePropTypes.map,

    patchSeason: PropTypes.func.isRequired,
    loadUserTeams: PropTypes.func.isRequired,
    createTeamSeasonRequest: PropTypes.func.isRequired,
    loadTeamSeasonRequests: PropTypes.func.isRequired,
    patchTeamSeasonRequest: PropTypes.func.isRequired,
    loadTeamSeasons: PropTypes.func.isRequired,
    patchTeamSeason: PropTypes.func.isRequired,
    loadNewsItems: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ season, newsItems }, prevProps) {
    if (!season) return;
    if (
      !prevProps ||
      season !== prevProps.season ||
      (!newsItems && prevProps.newsItems)
    ) {
      this.props.loadNewsItems({
        target: 'season',
        targetId: season.get('id'),
      });
    }
  }

  handleAcceptAll = (_, id) => {
    const { season } = this.props;
    if (!window.confirm('Accept all applicants?')) return;
    this.props.patchSeason(id, 'accept', () => {
      this.props.messagePush('accepted all applicants');
      if (!season) return;
      this.props.loadTeamSeasons(
        { seasonId: season.get('id') }, undefined, ['teamId'],
      );
      this.props.loadTeamSeasonRequests(
        { seasonId: season.get('id') }, undefined, ['teamId'],
      );
    });
  };

  handleEndSeason = (_, id) => {
    if (!window.confirm('End this season?')) return;
    this.props.patchSeason(id, 'end', () =>
      this.props.messagePush('ended a season'),
    );
  };

  render() {
    const {
      isLoading, myId, amAdmin, tournament: _tournament, season: _season,
      myParticipation, myApplication, leaderships,
      applications, teams, seasonPath, newsItems,
    } = this.props;

    if (!_tournament || !_season) {
      return isLoading ?
        <div className={this.cni()}><Loading /></div> :
          <NotFound />;
    }

    const tournament = _tournament.toJS();
    const season = _season.toJS();
    season.signupsOpenedAt =
      season.signupsOpenedAt && new Date(season.signupsOpenedAt);
    season.signupsClosedAt =
      season.signupsClosedAt && new Date(season.signupsClosedAt);

    let myTeamId;
    let canApply;
    if (
      myParticipation &&
      myApplication &&
      new Date(myParticipation.get('createdAt')) >
      new Date(myApplication.get('createdAt'))
    ) {
      myTeamId = myParticipation.get('teamId');
      canApply = myParticipation.get('leftAt') && (
        !myParticipation.get('isDone') ||
        myParticipation.get('kickedBy')
      );
    } else if (myApplication) {
      myTeamId = myApplication.get('teamId');
      canApply = myApplication.get('decision') === false;
    } else {
      canApply = true;
    }

    const myTeam = myTeamId && teams.get(myTeamId);
    const amLeader = leaderships && leaderships.size;
    const now = new Date();
    if (!amLeader) {
      canApply = false;
    } else if (!season.signupsOpenedAt || season.signupsOpenedAt > now) {
      canApply = false;
    } else if (season.signupsClosedAt && season.signupsClosedAt < now) {
      canApply = false;
    }

    const newsBlock = [];
    if (newsItems) {
      let size = 'L';
      for (const x of newsItems.toJS()) {
        const publishedAt = x.publishedAt && new Date(x.publishedAt);
        newsBlock.push(
          <Link
            key={x.id}
            className={this.cn({ d: 'tile', m: size })}
            to={`${seasonPath}/news/${x.id}`}
          >
            <div
              className={this.cn({ d: 'tileImage' })}
              style={{ backgroundImage: `url(${utils.https(x.picture)})` }}
            />
            <div className={this.cn({ d: 'tileInfo' })}>
              <div className={this.cn({ d: 'tileTitle' })}>
                {x.title}
              </div>
              <div className={this.cn({ d: 'tileDate' })}>{
                (
                  x.isDeleted &&
                  'DELETED'
                ) || (
                  publishedAt &&
                  publishedAt < now &&
                  moment(publishedAt).format('ll')
                ) || (
                  publishedAt &&
                  publishedAt > now &&
                  `SCHEDULED FOR ${moment(publishedAt).format('ll')}`
                ) || 'DRAFT'
              }</div>
            </div>
          </Link>
        );

        if (size === 'L') size = 'M';
        else if (size === 'M') size = 'S';
      }

      if (newsBlock.length === 1) {
        newsBlock.push(
          <div
            key="extra"
            className={this.cn({
              d: 'tile',
              m: ['M', 'none'],
            })}
          />
        );
      } else if ((newsBlock.length - 2) % 3 > 0) {
        for (let i = 0; i < (newsBlock.length - 2) % 3; i += 1) {
          newsBlock.push(
            <div
              key={`extra${i}`}
              className={this.cn({
                d: 'tile',
                m: ['S', 'none'],
              })}
            />
          );
        }
      }
    }

    return (
      <div className={this.cni()}>
        <div className={this.cn({ u: 'sectionSingle' })}>
          <div className={this.cn({ d: 'info' })}>
            <div className={this.cn({ u: 'sectionMargined' })}>
              <div className={this.cn({ d: 'infoInner' })}>
                <div className={this.cn({ d: 'infoHeading' })}>
                  INFO
                </div>
                <div className={this.cn({ d: 'infoItem' })}>
                  <div className={this.cn({ d: 'infoTitle' })}>
                    FORMAT
                  </div>
                  <div className={this.cn({ d: 'infoBody' })}>
                    {season.teamSize}v{season.teamSize}
                  </div>
                </div>
                <div className={this.cn({ d: 'infoItem' })}>
                  <div className={this.cn({ d: 'infoTitle' })}>
                    CAPACITY
                  </div>
                  <div className={this.cn({ d: 'infoBody' })}>
                    {season.capacity || '∞'} teams
                  </div>
                </div>
                <div className={this.cn({ d: 'infoItem' })}>
                  <div className={this.cn({ d: 'infoTitle' })}>
                    SIGNUPS OPEN
                  </div>
                  <div className={this.cn({ d: 'infoBody' })}>{
                    season.signupsOpenedAt ?
                      moment(season.signupsOpenedAt).format('lll') :
                      'TBA'
                  }</div>
                </div>
                <div className={this.cn({ d: 'infoItem' })}>
                  <div className={this.cn({ d: 'infoTitle' })}>
                    SIGNUPS CLOSE
                  </div>
                  <div className={this.cn({ d: 'infoBody' })}>{
                    season.signupsClosedAt ?
                      moment(season.signupsClosedAt).format('lll') :
                      'TBA'
                  }</div>
                </div>
                <div className={this.cn({ d: 'infoItem' })}>
                  <div className={this.cn({ d: 'infoTitle' })}>
                    DURATION
                  </div>
                  <div className={this.cn({ d: 'infoBody' })}>
                    {season.duration} days
                  </div>
                </div>
                <div className={this.cn({ d: 'infoItem' })}>
                  <div className={this.cn({ d: 'infoTitle' })}>
                    MAX TEAM SIZE
                  </div>
                  <div className={this.cn({ d: 'infoBody' })}>
                    {season.teamSizeMax || '∞'}
                  </div>
                </div>
                <div className={this.cn({ d: 'infoFiller' })} />
                {canApply && <Button
                  className={this.cn({ d: 'infoButton' })}
                  text="apply..." type="important"
                  href={`${seasonPath}/apply`}
                />}
                {!canApply && myTeam && <div
                  className={this.cn({ d: 'infoItem' })}
                >
                  <div className={this.cn({ d: 'infoTitle' })}>
                    REPRESENTING
                  </div>
                  <div className={this.cn({ d: 'infoBody', m: 'right' })}>
                    {myTeam.get('abbr')}
                  </div>
                </div>}
                {myId && !amLeader && <Button
                  className={this.cn({ d: 'infoButton' })}
                  text="create team to apply..." type="important"
                  href="/teams/new"
                />}
                {!myId && <Button
                  className={this.cn({ d: 'infoButton' })}
                  text="sign up to apply..." type="important"
                  href="/signup"
                />}
              </div>
            </div>
          </div>
          {(
            tournament.email ||
            tournament.discord ||
            tournament.twitter ||
            tournament.facebook ||
            tournament.twitch ||
            tournament.youtube
          ) && <div className={this.cn({ d: 'extra' })}>
            <div className={this.cn({ u: 'sectionMargined' })}>
              <div className={this.cn({ d: 'extraInner' })}>
                {tournament.email && <a
                  className={this.cn({ d: 'extraItem' })}
                  href={`mailto:${tournament.email}`}
                >
                  <img src={require('../assets/email.svg')} />
                  <div><div>{tournament.email}</div></div>
                </a>}
                {tournament.discord && <a
                  className={this.cn({ d: 'extraItem' })}
                  href={`//discord.gg/${tournament.discord}`}
                >
                  <img src={require('../assets/discord.svg')} />
                  <div><div>community chat</div></div>
                </a>}
                {tournament.twitter && <a
                  className={this.cn({ d: 'extraItem' })}
                  href={`//twitter.com/${tournament.twitter}`}
                >
                  <img src={require('../assets/twitter.svg')} />
                  <div><div>@{tournament.twitter}</div></div>
                </a>}
                {tournament.facebook && <a
                  className={this.cn({ d: 'extraItem' })}
                  href={`//facebook.com/${tournament.facebook}`}
                >
                  <img src={require('../assets/facebook.svg')} />
                  <div><div>{tournament.facebook}</div></div>
                </a>}
                {tournament.twitch && <a
                  className={this.cn({ d: 'extraItem' })}
                  href={`//twitch.tv/${tournament.twitch}`}
                >
                  <img src={require('../assets/twitch.svg')} />
                  <div><div>{tournament.twitch}</div></div>
                </a>}
                {tournament.youtube && <a
                  className={this.cn({ d: 'extraItem' })}
                  href={`//youtube.com/c/${tournament.youtube}`}
                >
                  <img src={require('../assets/youtube.svg')} />
                  <div><div>{tournament.youtube}</div></div>
                </a>}
              </div>
            </div>
          </div>}
          <div className={this.cn({ d: 'about' })}>
            <div className={this.cn({ u: 'sectionMargined', d: 'aboutInner' })}>
              <Markdown
                className={this.cn({ d: 'aboutBody' })}
                source={`${tournament.description}\n\n${season.description}`}
              />
              <div className={this.cn({ d: 'aboutFill' })} />
              <div className={this.cn({ d: 'aboutLogo' })}>
                <div
                  className={this.cn({ d: 'aboutLogoBackground' })}
                  style={{
                    backgroundImage: `url(${utils.https(tournament.logo)})`,
                  }}
                />
                <div className={this.cn({ d: 'aboutLogoTag' })}>
                  {!tournament.logoHasText && tournament.abbr}
                </div>
              </div>
            </div>
          </div>
          {tournament.twitchLive && <div className={this.cn({ d: 'live' })}>
            <div className={this.cn({ d: 'liveInner' })}>
              <div className={this.cn({ d: 'liveVideoOuter2' })}>
                <div className={this.cn({ d: 'liveVideoOuter' })}>
                  <iframe
                    className={this.cn({ d: 'liveVideo' })}
                    src={`https://player.twitch.tv/?channel=${tournament.twitchLive}`}
                    frameBorder="0"
                    scrolling="no"
                    allowFullScreen="true"
                  />
                </div>
              </div>
              <div className={this.cn({ d: 'liveChatOuter' })}>
                <iframe
                  className={this.cn({ d: 'liveChat' })}
                  src={`https://www.twitch.tv/${tournament.twitchLive}/chat`}
                  frameBorder="0"
                  scrolling="no"
                />
              </div>
            </div>
          </div>}
          <div className={this.cn({ d: 'news' })}>
            <div className={this.cn({ u: 'sectionMargined' })}>
              <div className={this.cn({ d: 'newsInner' })}>
                {newsBlock.length ? newsBlock : 'NO NEWS'}
              </div>
            </div>
          </div>
          {amAdmin && <div className={this.cn({ d: 'admin' })}>
            <div className={this.cn({ u: 'sectionMargined' })}>
              <div className={this.cn({ d: 'adminInner' })}>
                <div className={this.cn({ d: 'adminHeading' })}>
                  ADMIN
                </div>
                <Button
                  text="add news..." type="important" size="large"
                  href={`${this.props.seasonPath}/news/new`}
                />
                <Button
                  text="tournament settings..." type="important" size="large"
                  href={`${this.props.tournamentPath}/_/settings`}
                />
                <Button
                  text="season settings..." type="important" size="large"
                  href={`${this.props.seasonPath}/_/settings`}
                />
                {!season.endedAt &&
                season.signupsClosedAt &&
                season.signupsClosedAt < now &&
                applications &&
                applications.filter(
                  x => x.get('decision') === null
                ).size > 0 && <Button
                  text="accept all" type="important" size="large"
                  meta={season.id}
                  onClick={this.handleAcceptAll}
                />}
                {!season.endedAt && <Button
                  text="end season" type="important" size="large"
                  meta={season.id}
                  onClick={this.handleEndSeason}
                />}
              </div>
            </div>
          </div>}
        </div>
      </div>
    );
  }
}
