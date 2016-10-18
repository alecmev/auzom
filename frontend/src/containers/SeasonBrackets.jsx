import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import * as utils from '../utils';
import Component from '../utils/Component';

const bracketsPathSelector = (_, { seasonPath }) =>
  `${seasonPath}/brackets`;

const seasonStagesSelector = createSelector(
  selectors.stages,
  (_, props) => props.season && props.season.get('id'),
  (x, seasonId) =>
    seasonId && x.toList()
      .filter(y => y.get('seasonId') === seasonId)
      .sortBy(y => y.get('startedAt')),
);

const stageSlugSelector = (_, props) => props.params.stageSlug;
const stagePathSelector = (_, { seasonPath, params }) =>
  `${seasonPath}/brackets/${params.stageSlug}`;
const stageSelector = createSelector(
  seasonStagesSelector,
  stageSlugSelector,
  (x, slug) => x && slug && x.find(y => y.get('slug') === slug),
);

const stageBracketsSelector = createSelector(
  selectors.brackets,
  createSelector(stageSelector, x => x && x.get('id')),
  (x, stageId) =>
    stageId && x.toList()
      .filter(y => y.get('stageId') === stageId)
      .sortBy(y => y.get('order')),
);

const bracketSlugSelector = (_, props) => props.params.bracketSlug;
const bracketPathSelector = (_, { seasonPath, params }) =>
  `${seasonPath}/brackets/${params.stageSlug}/${params.bracketSlug}`;
const bracketSelector = createSelector(
  stageBracketsSelector,
  bracketSlugSelector,
  (x, slug) => x && slug && x.find(y => y.get('slug') === slug),
);

const bracketIdSelector = createSelector(
  bracketSelector, x => x && x.get('id'),
);

const bracket2RoundsSelector = createSelector(
  selectors.bracketRounds,
  bracketIdSelector,
  (x, bracketId) =>
    bracketId && x.toList().filter(y => y.get('bracketId') === bracketId),
);

const bracketMatchesSelector = createSelector(
  selectors.matches,
  bracketIdSelector,
  (x, bracketId) =>
    bracketId && x.toList().filter(y => y.get('bracketId') === bracketId),
);

const bracket2StandingsSelector = createSelector(
  selectors.bracketStandings,
  bracketIdSelector,
  (x, id) => id && x.getIn([id, 'standings']),
);

@connect(createStructuredSelector({
  amAdmin: selectors.amAdmin,

  bracketsPath: bracketsPathSelector,
  seasonStages: seasonStagesSelector,
  stageSlug: stageSlugSelector,
  stagePath: stagePathSelector,
  stage: stageSelector,
  stageBrackets: stageBracketsSelector,
  bracketSlug: bracketSlugSelector,
  bracketPath: bracketPathSelector,
  bracket: bracketSelector,
  bracketId: bracketIdSelector,
  bracket2Rounds: bracket2RoundsSelector,
  bracketMatches: bracketMatchesSelector,
  bracket2Standings: bracket2StandingsSelector,
}), actions)
export default class SeasonBrackets extends Component {
  static propTypes = {
    children: PropTypes.element,
    game: ImmutablePropTypes.map,
    seasonPath: PropTypes.string,
    season: ImmutablePropTypes.map,

    amAdmin: PropTypes.bool,

    bracketsPath: PropTypes.string.isRequired,
    seasonStages: ImmutablePropTypes.list,
    stageSlug: PropTypes.string,
    stagePath: PropTypes.string.isRequired,
    stage: ImmutablePropTypes.map,
    stageBrackets: ImmutablePropTypes.list,
    bracketSlug: PropTypes.string,
    bracketPath: PropTypes.string.isRequired,
    bracket: ImmutablePropTypes.map,
    bracket2Rounds: ImmutablePropTypes.list,
    bracketMatches: ImmutablePropTypes.list,
    bracket2Standings: ImmutablePropTypes.list,

    loadStages: PropTypes.func.isRequired,
    loadBrackets: PropTypes.func.isRequired,
    loadBracketRounds: PropTypes.func.isRequired,
    loadMatches: PropTypes.func.isRequired,
    loadTeam: PropTypes.func.isRequired,
    loadBracketStandings: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
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
    tournament, season, bracketsPath, seasonStages, stagePath, stage,
    stageBrackets, bracket, bracket2Rounds, bracketMatches, bracket2Standings,
    ...nextProps
  }, prevProps) {
    if (!season) return;
    if (
      !prevProps ||
      season !== prevProps.season ||
      (!seasonStages && prevProps.seasonStages)
    ) {
      this.props.loadStages({ seasonId: season.get('id') });
    }

    if (!stage) {
      if (
        !seasonStages || !seasonStages.size ||
        React.Children.count(nextProps.children) // for stuff like settings
      ) return;

      const now = new Date();
      let currentStage;
      seasonStages.forEach((x) => {
        if (new Date(x.get('startedAt')) > now) {
          if (!currentStage) currentStage = x;
          return false;
        }

        currentStage = x;
        return true;
      });

      this.context.router.replace(
        `${bracketsPath}/${currentStage.get('slug')}`,
      );
      return;
    }

    if (
      !prevProps ||
      stage !== prevProps.stage ||
      (!stageBrackets && prevProps.stageBrackets)
    ) {
      this.props.loadBrackets({ stageId: stage.get('id') });
    }

    if (!bracket) {
      if (
        !stageBrackets || !stageBrackets.size ||
        React.Children.count(nextProps.children) // for stuff like settings
      ) return;

      // TODO: deduct my team's bracket
      this.context.router.replace(
        `${stagePath}/${stageBrackets.first().get('slug')}`,
      );
      return;
    }

    const force = !prevProps || bracket !== prevProps.bracket;
    const bracketId = bracket.get('id');
    if (force || (!bracket2Rounds && prevProps.bracket2Rounds)) {
      this.props.loadBracketRounds({ bracketId });
    }
    if (force || (!bracketMatches && prevProps.bracketMatches)) {
      this.props.loadMatches({ bracketId });
    }
    if (
      force ||
      (!bracket2Standings && prevProps.bracket2Standings) ||
      bracketMatches !== prevProps.bracketMatches
    ) {
      this.props.loadBracketStandings(bracketId);
    }

    if (!bracket2Standings) return;
    if (!prevProps || bracket2Standings !== prevProps.bracket2Standings) {
      bracket2Standings.forEach(x => this.props.loadTeam(x.get('teamId')));
    }
  }

  render() {
    const {
      seasonStages, amAdmin, bracketsPath, stage, stageBrackets, bracket,
      bracketMatches,
    } = this.props;

    const childProps = utils.cloneProps(this.props);
    const children = React.Children.map(
      this.props.children, x => React.cloneElement(x, childProps),
    );

    const type = bracket && bracket.get('type');
    const isSwiss = type && (
      type === 'bcl-sc16-swiss' ||
      type === 'ace-pre-swiss'
    );

    return (
      <div className={this.cni()}>
        <div className={this.cn({ u: 'sectionSingle', d: 'menu' })}>
          <div className={this.cn({ u: 'sectionMargined' })}>
            <div className={this.cn({ d: 'menuInner' })}>
              <div className={this.cn({ d: 'menuItem', m: 'title' })}>
                stage
              </div>
              <div className={this.cn({ d: 'menuDivider' })} />
              {seasonStages && seasonStages.map(x =>
                <Link
                  key={x.get('id')}
                  to={`${bracketsPath}/${x.get('slug')}`}
                  className={this.cn({ d: 'menuItem' })}
                  activeClassName="is-active"
                >
                  {x.get('name')}
                </Link>
              ).toJS()}
              {amAdmin && <Link
                to={`${bracketsPath}/_/new-stage`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                + new stage
              </Link>}
            </div>
          </div>
          {stage && stageBrackets && <div className={this.cn({
            u: 'sectionMargined',
          })}>
            <div className={this.cn({ d: 'menuInner' })}>
              <div className={this.cn({ d: 'menuItem', m: 'title' })}>
                bracket
              </div>
              <div className={this.cn({ d: 'menuDivider' })} />
              {stageBrackets.map(x =>
                <Link
                  key={x.get('id')}
                  to={`${bracketsPath}/${stage.get('slug')}/${x.get('slug')}`}
                  className={this.cn({ d: 'menuItem' })}
                  activeClassName="is-active"
                >
                  {x.get('name')}
                </Link>
              ).toJS()}
              {amAdmin && <Link
                to={`${bracketsPath}/${stage.get('slug')}/_/new-bracket`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                + new bracket
              </Link>}
            </div>
          </div>}
          {amAdmin && stage && <div className={this.cn({
            u: 'sectionMargined',
          })}>
            <div className={this.cn({ d: 'menuInner' })}>
              <div className={this.cn({ d: 'menuItem', m: 'title' })}>
                admin
              </div>
              <div className={this.cn({ d: 'menuDivider' })} />
              <Link
                to={`${bracketsPath}/${stage.get('slug')}/_/settings`}
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                stage settings
              </Link>
              {bracket && <Link
                to={
                  `${bracketsPath}/${stage.get('slug')}/` +
                  `${bracket.get('slug')}/_/settings`
                }
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                bracket settings
              </Link>}
              {bracket && !isSwiss &&
              bracketMatches && !!bracketMatches.size && (() => {
                const match = bracketMatches.first().toJS();
                if (match.teamX || match.teamY) return null;
                return (
                  <Link
                    to={
                      `${bracketsPath}/${stage.get('slug')}/` +
                      `${bracket.get('slug')}/_/prepare`
                    }
                    className={this.cn({ d: 'menuItem' })}
                    activeClassName="is-active"
                  >
                    prepare the bracket
                  </Link>
                );
              })()}
              {bracket && isSwiss && <Link
                to={
                  `${bracketsPath}/${stage.get('slug')}/` +
                  `${bracket.get('slug')}/_/new-swiss-round`
                }
                className={this.cn({ d: 'menuItem' })}
                activeClassName="is-active"
              >
                + new swiss round
              </Link>}
            </div>
          </div>}
        </div>
        {children}
      </div>
    );
  }
}
