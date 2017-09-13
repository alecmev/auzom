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

import Button from '../components/Button';
import Img from '../components/Img';

const teamIdSelector = (_, props) => props.params.teamId;

const teamSelector = createSelector(
  selectors.teams,
  teamIdSelector,
  utils.get,
);

const teamSeasonsSelector = createSelector(
  (_, props) => props.participants,
  teamIdSelector,
  (x, teamId) => x && x
    .filter(y => y.get('teamId') === teamId)
    .sort((a, b) => b.get('createdAt').localeCompare(a.get('createdAt'))),
);

const teamSeasonRequestsSelector = createSelector(
  (_, props) => props.applications,
  teamIdSelector,
  (x, teamId) => x && x
    .filter(y => y.get('teamId') === teamId)
    .sort((a, b) => b.get('createdAt').localeCompare(a.get('createdAt'))),
);

@connect(createStructuredSelector({
  amAdmin: selectors.amAdmin,
  teamId: teamIdSelector,
  team: teamSelector,
  teamSeasons: teamSeasonsSelector,
  teamSeasonRequests: teamSeasonRequestsSelector,
}), actions)
export default class SeasonParticipant extends Component {
  static propTypes = {
    location: PropTypes.object, // eslint-disable-line
    season: ImmutablePropTypes.map,

    amAdmin: PropTypes.bool,
    teamId: PropTypes.string,
    team: ImmutablePropTypes.map,
    teamSeasons: ImmutablePropTypes.list,
    teamSeasonRequests: ImmutablePropTypes.list,
    leaderships: ImmutablePropTypes.list,

    loadTeamSeasons: PropTypes.func.isRequired,
    patchTeamSeason: PropTypes.func.isRequired,
    patchTeamSeasonRequest: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  // no need in loading, it's already done by Tournament

  handleDone = (_, id) => {
    if (!window.confirm('Mark as done?')) return;
    this.props.patchTeamSeason(id, 'done', () =>
      this.props.messagePush('marked a team as done'),
    );
  };

  handleLeave = (_, id) => {
    if (!window.confirm('Mark as left?')) return;
    this.props.patchTeamSeason(id, 'leave', () =>
      this.props.messagePush('marked a team as left'),
    );
  };

  handleKick = (_, id) => {
    if (!window.confirm('Kick?')) return;
    this.props.patchTeamSeason(id, 'kick', () =>
      this.props.messagePush('kicked a team from a season'),
    );
  };

  handleAccept = (_, id) => {
    if (!window.confirm('Accept?')) return;
    this.props.patchTeamSeasonRequest(id, 'yes', (x) => {
      x.decision && this.props.loadTeamSeasons({
        teamId: x.teamId,
        seasonId: x.seasonId,
      });
      this.props.messagePush('accepted an application to a season');
    });
  };

  handleDecline = (_, id) => {
    if (!window.confirm('Decline?')) return;
    this.props.patchTeamSeasonRequest(id, 'no', () =>
      this.props.messagePush('declined an application to a season'),
    );
  };

  handleCancel = (_, id) => {
    if (!window.confirm('Cancel?')) return;
    this.props.patchTeamSeasonRequest(id, 'cancel', () =>
      this.props.messagePush('cancelled an application to a season'),
    );
  };

  render() {
    const {
      amAdmin, season, teamId, team, teamSeasons, teamSeasonRequests,
      leaderships,
    } = this.props;
    if (!season || !team) return null;
    const canAccept = new Date(season.get('signupsClosedAt')) < new Date();
    const amLeader = leaderships && leaderships.find(
      y => y.get('teamId') === teamId,
    );
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined' })}>
          <Link
            className={this.cn({ d: 'heading' })}
            to={`/teams/${teamId}`}
          >
            <div className={this.cn({ d: 'logoOuter' })}>
              {team.get('logo') && <Img
                className={this.cn({ d: 'logo' })}
                src={team.get('logo')}
              />}
            </div>
            <div className={this.cn({ d: 'name' })}>
              <div className={this.cn({ d: 'nameTag' })}>
                {team.get('abbr')}
              </div>
              <div className={this.cn({ d: 'nameFull' })}>
                {team.get('name')}
              </div>
            </div>
          </Link>
          {teamSeasons && teamSeasons.toJS().map((x, i) => (
            <div key={x.id} className={this.cn({ d: 'info' })}>
              <div className={this.cn({ d: 'infoHeading' })}>
                PARTICIPATION {teamSeasons.size > 1 && (teamSeasons.size - i)}
              </div>
              <div className={this.cn({ d: 'infoItem' })}>
                <div className={this.cn({ d: 'infoTitle' })}>
                  ACCEPTED
                </div>
                <div className={this.cn({ d: 'infoBody' })}>
                  {moment(x.createdAt).format('lll')}
                </div>
              </div>
              {x.leftAt && <div
                className={this.cn({ d: 'infoItem' })}
              >
                <div className={this.cn({ d: 'infoTitle' })}>
                  LEFT
                </div>
                <div className={this.cn({ d: 'infoBody' })}>
                  {moment(x.leftAt).format('lll')}
                </div>
              </div>}
              {x.leftAt && <div
                className={this.cn({ d: 'infoItem' })}
              >
                <div className={this.cn({ d: 'infoTitle' })}>
                  LEAVE REASON
                </div>
                <div className={this.cn({
                  d: 'infoBody',
                  m: !x.isDone && 'red',
                })}>
                  {(x.isDone && 'FINISHED') ||
                  (x.kickedBy && 'KICKED') ||
                  'QUIT'}
                </div>
              </div>}
              <div className={this.cn({ d: 'infoFiller' })} />
              {amAdmin && !x.leftAt &&
                <div className={this.cn({ d: 'admin' })}>
                  <Button
                    text="done"
                    type="important"
                    size="large"
                    meta={x.id}
                    onClick={this.handleDone}
                  />
                  <Button
                    text="leave"
                    type="important"
                    size="large"
                    meta={x.id}
                    onClick={this.handleLeave}
                  />
                  <Button
                    text="kick"
                    type="important"
                    size="large"
                    meta={x.id}
                    onClick={this.handleKick}
                  />
                </div>
              }
            </div>
          ))}
          {teamSeasonRequests && teamSeasonRequests.toJS().map((x, i) => (
            <div key={x.id} className={this.cn({ d: 'info' })}>
              <div className={this.cn({ d: 'infoHeading' })}>
                APPLICATION {
                  teamSeasonRequests.size > 1 && (teamSeasonRequests.size - i)
                }
              </div>
              <div className={this.cn({ d: 'infoItem' })}>
                <div className={this.cn({ d: 'infoTitle' })}>
                  SUBMITTED
                </div>
                <div className={this.cn({ d: 'infoBody' })}>
                  {moment(x.createdAt).format('lll')}
                </div>
              </div>
              {x.cancelledBy && <div
                className={this.cn({ d: 'infoItem' })}
              >
                <div className={this.cn({ d: 'infoTitle' })}>
                  CANCELLED
                </div>
                <div className={this.cn({ d: 'infoBody' })}>
                  {moment(x.decidedAt).format('lll')}
                </div>
              </div>}
              {!x.cancelledBy && x.decidedAt && <div
                className={this.cn({ d: 'infoItem' })}
              >
                <div className={this.cn({ d: 'infoTitle' })}>
                  REVIEWED
                </div>
                <div className={this.cn({ d: 'infoBody' })}>
                  {moment(x.decidedAt).format('lll')}
                </div>
              </div>}
              {!x.cancelledBy && x.decidedAt && <div
                className={this.cn({ d: 'infoItem' })}
              >
                <div className={this.cn({ d: 'infoTitle' })}>
                  DECISION
                </div>
                <div className={this.cn({
                  d: 'infoBody',
                  m: x.decision ? 'green' : 'red',
                })}>
                  {x.decision ? 'ACCEPTED' : 'DECLINED'}
                </div>
              </div>}
              <div className={this.cn({ d: 'infoFiller' })} />
              {x.decision === null &&
                <div className={this.cn({ d: 'admin' })}>
                  {amAdmin && canAccept && <Button
                    text="accept"
                    type="important"
                    size="large"
                    meta={x.id}
                    onClick={this.handleAccept}
                  />}
                  {amAdmin && <Button
                    text="decline"
                    type="important"
                    size="large"
                    meta={x.id}
                    onClick={this.handleDecline}
                  />}
                  {amLeader && <Button
                    text="cancel"
                    type="important"
                    size="large"
                    meta={x.id}
                    onClick={this.handleCancel}
                  />}
                </div>
              }
            </div>
          ))}
        </div>
      </div>
    );
  }
}
