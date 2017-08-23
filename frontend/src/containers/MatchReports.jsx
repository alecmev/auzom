import moment from 'moment';
import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';

import MatchReportsItem from './MatchReportsItem';

const matchReportsSelector = createSelector(
  selectors.matchReports,
  (_, props) => props.matchId,
  (x, matchId) =>
    matchId && x.toList().filter(y => y.get('matchId') === matchId),
);

@connect(createStructuredSelector({
  matchReports: matchReportsSelector,
}), actions)
export default class MatchReports extends Component {
  static propTypes = {
    matchId: PropTypes.string,
    match: ImmutablePropTypes.map,
    matchLeadership: ImmutablePropTypes.list,
    linkPrefix: PropTypes.string,

    matchReports: ImmutablePropTypes.list,

    loadMatchReports: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ matchId, matchReports }, prevProps) {
    if (!matchId) return;
    if (
      !prevProps || matchId !== prevProps.matchId ||
      (!matchReports && prevProps.matchReports)
    ) {
      this.props.loadMatchReports(
        { matchId }, undefined, ['createdBy', 'agreedUponBy'],
      );
    }
  }

  render() {
    const {
      matchLeadership, matchReports, match, matchId, linkPrefix,
    } = this.props;

    if (!match || !matchLeadership) return null;

    const myTeamId = (matchLeadership.size === 1) ?
      matchLeadership.first() : null;
    const currentReport = match.get('matchReportId');
    const rd = match.get('reportingClosedAt');
    const rdMoment = rd && moment(rd);

    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined' })}>
          <div>
            {rd && rdMoment.isBefore() ?
              'Reporting closed at:' :
              'Report deadline:'
            }
            <span className={this.cn({ d: 'deadline' })}>
              {rd ? rdMoment.format('lll') : 'stage end'}
            </span>
          </div>
          {matchReports && (!matchReports.size ?
            'No reports' :
            <ul>{matchReports.map(x =>
              <MatchReportsItem
                key={x.get('id')}
                matchReport={x}
                publishedReportId={currentReport}
                myTeamId={myTeamId}
                matchId={matchId}
                linkPrefix={linkPrefix}
              />
            )}</ul>
          )}
        </div>
      </div>
    );
  }
}
