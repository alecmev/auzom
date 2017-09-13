import { fromJS } from 'immutable';
import { groupBy } from 'lodash';
import moment from 'moment';
import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Link, IndexLink } from 'react-router';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import * as utils from '../utils';
import Component from '../utils/Component';

const matchIdSelector = (_, props) => props.params.matchId;
const matchPathSelector = (_, { bracketPath, params }) =>
  `${bracketPath}/matches/${params.matchId}`;
const matchSelector = createSelector(
  selectors.matches,
  matchIdSelector,
  utils.get,
);

const bracketRoundSelector = createSelector(
  (_, props) => props.bracket2Rounds,
  matchSelector,
  (bracket2Rounds, match) => bracket2Rounds && match && bracket2Rounds.find(
    x =>
      x.get('bracketId') === match.get('bracketId') &&
      x.get('number') === match.get('bracketRound'),
  ),
);

const matchLeadershipSelector = createSelector(
  selectors.matchLeaderships,
  matchIdSelector,
  (x, id) => id && x.getIn([id, 'leadership']),
);

const matchMapsSelector = createSelector(
  selectors.matchMaps,
  matchIdSelector,
  (x, id) => id && x.toList().filter(
    y => y.get('matchId') === id && y.get('discardedAt') === null,
  ).sortBy(y => +y.get('id')), // TODO: hmm
);

const matchReportIdSelector = createSelector(
  matchSelector, x => x && x.get('matchReportId'),
);

const matchRoundsSelector = createSelector(
  selectors.matchRounds,
  matchReportIdSelector,
  (x, id) => id && x.toList().filter(
    y => y.get('matchReportId') === id,
  ).sortBy(y => +y.get('id')), // TODO: hmm
);

const matchPenaltiesSelector = createSelector(
  selectors.matchPenalties,
  matchReportIdSelector,
  (x, id) => id && x.toList().filter(y => y.get('matchReportId') === id),
);

const matchWidePenaltiesSelector = createSelector(
  matchPenaltiesSelector,
  x => x && x.filter(y => y.get('matchRoundId') === null),
);

const roundWidePenaltiesSelector = createSelector(
  matchPenaltiesSelector,
  x => x && fromJS(groupBy(
    x.filter(y => y.get('matchRoundId') !== null).toJS(), 'matchRoundId',
  )),
);

@connect(createStructuredSelector({
  amAdmin: selectors.amAdmin,
  teams: selectors.teams,
  gameMaps: selectors.gameMaps,
  matchReports: selectors.matchReports,

  matchId: matchIdSelector,
  matchPath: matchPathSelector,
  match: matchSelector,
  bracketRound: bracketRoundSelector,
  matchLeadership: matchLeadershipSelector,
  matchMaps: matchMapsSelector,
  matchReportId: matchReportIdSelector,
  matchRounds: matchRoundsSelector,
  matchWidePenalties: matchWidePenaltiesSelector,
  roundWidePenalties: roundWidePenaltiesSelector,
}), actions)
export default class Match extends Component {
  static propTypes = {
    children: PropTypes.element,
    game: ImmutablePropTypes.map,
    seasonPath: PropTypes.string,
    bracketPath: PropTypes.string,

    amAdmin: PropTypes.bool,
    teams: ImmutablePropTypes.map,
    gameMaps: ImmutablePropTypes.map,
    matchReports: ImmutablePropTypes.map,

    matchId: PropTypes.string,
    matchPath: PropTypes.string,
    match: ImmutablePropTypes.map,
    matchLeadership: ImmutablePropTypes.list,
    matchMaps: ImmutablePropTypes.list,
    matchReportId: PropTypes.string,
    matchRounds: ImmutablePropTypes.list,
    matchWidePenalties: ImmutablePropTypes.list,
    roundWidePenalties: ImmutablePropTypes.map,

    loadMatchLeadership: PropTypes.func.isRequired,
    loadMatchMaps: PropTypes.func.isRequired,
    loadMatchReport: PropTypes.func.isRequired,
    loadMatchRounds: PropTypes.func.isRequired,
    loadMatchPenalties: PropTypes.func.isRequired,
    patchMatch: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ matchId, matchLeadership, matchMaps, matchReportId }, prevProps) {
    if (!matchId) return;
    const force = !prevProps || matchId !== prevProps.matchId;
    if (force || (!matchLeadership && prevProps.matchLeadership)) {
      this.props.loadMatchLeadership(matchId);
    }
    if (force || (!matchMaps && prevProps.matchMaps)) {
      this.props.loadMatchMaps(
        { matchId, discardedAt: '%00' }, undefined, ['gameMapId'],
      );
    }

    if (!matchReportId) return;
    if (!prevProps || matchReportId !== prevProps.matchReportId) {
      this.props.loadMatchReport(matchReportId);
      this.props.loadMatchRounds({ matchReportId });
      this.props.loadMatchPenalties({ matchReportId });
    }
  }

  handleResetMaps = () => {
    if (!this.props.match) return;
    if (!window.confirm('Reset this match\'s maps?')) return;
    this.props.patchMatch(
      this.props.match.get('id'), { action: 'reset-maps' },
      () => { location.href = location.href; },
    );
  };

  render() {
    const {
      amAdmin, teams, gameMaps, match: _match, matchLeadership, matchMaps,
      matchReports, matchReportId, matchPath, seasonPath,
    } = this.props;

    if (!_match) return null; // TODO: loading
    const match = _match.toJS();

    const teamX = match.teamX && teams.get(match.teamX);
    const teamY = match.teamY && teams.get(match.teamY);

    let winner;
    if (match.scoreX > match.scoreY) winner = 'x';
    else if (match.scoreX < match.scoreY) winner = 'y';
    else if (match.scoreX !== null) winner = 'draw';

    const getMod = (side) => {
      if (!winner) return winner;
      if (winner === 'draw') return winner;
      if (winner === side) return 'win';
      return 'loss';
    };

    const [modX, modY] = [getMod('x'), getMod('y')];

    let amLeader = false;
    if (matchLeadership && matchLeadership.size > 0) amLeader = true;

    const childProps = utils.cloneProps(this.props);
    childProps.match = _match;
    childProps.matchReport = matchReports.get(matchReportId);
    const children = React.Children.count(this.props.children) ?
      React.Children.map(
        this.props.children, x => React.cloneElement(x, childProps),
      ) : null;

    const when = moment(match.startedAt);

    return (
      <div className={this.cni()}>
        <div className={this.cn({ d: 'main', u: 'sectionSingle' })}>
          <div className={this.cn({ d: 'team', m: modX })}>
            <Link
              className={this.cn({ d: 'name' })}
              to={`${seasonPath}/participants/${match.teamX}`}
            >
              <div className={this.cn({ d: 'tag' })}>
                {teamX && teamX.get('abbr')}
              </div>
              {teamX && teamX.get('name')}
            </Link>
            <div className={this.cn({ d: 'logoPlaceholder' })} />
            {winner &&
              <div className={this.cn({ d: 'score', m: modX })}>
                {match.scoreX}
              </div>
            }
          </div>
          {!winner &&
            <div className={this.cn({ d: 'vs' })}>VS</div>
          }
          <div className={this.cn({ d: 'team', m: modY })}>
            {winner &&
              <div className={this.cn({ d: 'score', m: modY })}>
                {match.scoreY}
              </div>
            }
            <div className={this.cn({ d: 'logoPlaceholder' })} />
            <Link
              className={this.cn({ d: 'name' })}
              to={`${seasonPath}/participants/${match.teamY}`}
            >
              <div className={this.cn({ d: 'tag' })}>
                {teamY && teamY.get('abbr')}
              </div>
              {teamY && teamY.get('name')}
            </Link>
          </div>
        </div>
        <div className={this.cn({ u: 'sectionSingle', d: 'menu' })}>
          <div className={this.cn({ u: 'sectionMargined' })}>
            <div className={this.cn({ d: 'menuInner' })}>
              <div className={this.cn({ d: 'menuItem', m: 'title' })}>
                match
              </div>
              <div className={this.cn({ d: 'menuDivider' })} />
              <IndexLink
                to={matchPath}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                comments
              </IndexLink>
              {match.matchReportId && <Link
                to={`${matchPath}/details`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                details
              </Link>}
              <Link
                to={`${matchPath}/maps`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                maps
              </Link>
            </div>
          </div>
          {(amAdmin || amLeader) && <div className={this.cn({
            u: 'sectionMargined',
          })}>
            <div className={this.cn({ d: 'menuInner' })}>
              <div className={this.cn({ d: 'menuItem', m: 'title' })}>
                leader
              </div>
              <div className={this.cn({ d: 'menuDivider' })} />
              <IndexLink
                to={`${matchPath}/reports`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                reports
              </IndexLink>
              <Link
                to={`${matchPath}/reports/new`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                + new report
              </Link>
              {amLeader && <IndexLink
                to={`${matchPath}/attention-requests`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                attention requests
              </IndexLink>}
              {amLeader && <Link
                to={`${matchPath}/attention-requests/new`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                + request admin attention
              </Link>}
            </div>
          </div>}
          {amAdmin && <div className={this.cn({
            u: 'sectionMargined',
          })}>
            <div className={this.cn({ d: 'menuInner' })}>
              <div className={this.cn({ d: 'menuItem', m: 'title' })}>
                admin
              </div>
              <div className={this.cn({ d: 'menuDivider' })} />
              <Link
                to={`${matchPath}/_/settings`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                settings
              </Link>
              <Link
                to={`${matchPath}/_/prepare`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                prepare
              </Link>
              <button
                className={this.cn({ d: 'menuItem' })}
                onClick={this.handleResetMaps}
              >
                reset maps
              </button>
            </div>
          </div>}
        </div>
        <div className={this.cn({
          u: 'sectionSingle',
          d: 'context',
        })}>
          <div className={this.cn({ d: 'contextWhen' })}>
            {when.format('D MMMM YYYY')}
            <span className={this.cn({ d: 'contextWhenAt' })}> at </span>
            {when.format('HH:mm')}
          </div>
          <div className={this.cn({ d: 'contextWhere' })}>
            {matchMaps && matchMaps.filter(x => !x.get('isBan')).toJS().map(x =>
              <span key={x.id}>{gameMaps.getIn([x.gameMapId, 'name'])}</span>,
            )}
          </div>
        </div>
        {children}
      </div>
    );
  }
}
