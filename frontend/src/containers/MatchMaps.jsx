import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import { beep } from '../utils';
import Component from '../utils/Component';

const bracketMapsSelector = createSelector(
  selectors.bracketMaps,
  (_, props) => props.bracketId,
  (x, bracketId) =>
    bracketId && x.toList().filter(y =>
      y.get('bracketId') === bracketId && y.get('isEnabled'),
    ),
);

@connect(createStructuredSelector({
  gameMaps: selectors.gameMaps,
  bracketMaps: bracketMapsSelector,
}), actions)
export default class MatchMaps extends Component {
  static propTypes = {
    bracket: ImmutablePropTypes.map,
    bracketId: PropTypes.string,
    bracketRound: ImmutablePropTypes.map,
    match: ImmutablePropTypes.map,
    matchId: PropTypes.string,
    matchLeadership: ImmutablePropTypes.list,
    matchMaps: ImmutablePropTypes.list,

    gameMaps: ImmutablePropTypes.map,
    bracketMaps: ImmutablePropTypes.list,

    loadBracketMaps: PropTypes.func.isRequired,
    loadMatch: PropTypes.func.isRequired,
    loadMatchMaps: PropTypes.func.isRequired,
    patchMatch: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  state = {
    skipBeep: false,
  };

  componentWillMount() {
    this.load(this.props);
    const { bracket, match } = this.props;
    if (!bracket || !match || (
        bracket.get('mapVetoProcedure') !== '' &&
        !match.get('areMapsReady')
    )) {
      this.setState({ interval: setInterval(this.refresh, 10000) });
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.matchMaps !== this.props.matchMaps) {
      if (!this.state.skipBeep) beep();
      else this.setState({ skipBeep: false });
    }

    this.load(nextProps, this.props);
    const { bracket, match } = nextProps;
    if (bracket && match && (
        bracket.get('mapVetoProcedure') === '' ||
        match.get('areMapsReady')
    )) {
      clearInterval(this.state.interval);
    }
  }

  componentWillUnmount() {
    clearInterval(this.state.interval);
  }

  load({ bracketId, bracketMaps }, prevProps) {
    if (!bracketId) return;
    if (
      !prevProps || bracketId !== prevProps.bracketId ||
      (!bracketMaps && prevProps.bracketMaps)
    ) {
      this.props.loadBracketMaps({ bracketId }, null, ['gameMapId']);
    }
  }

  refresh = () => {
    const { matchId } = this.props;
    this.props.loadMatch(matchId);
    this.props.loadMatchMaps({ matchId, discardedAt: '%00' });
  };

  pick = (e) => {
    this.setState({ skipBeep: true });
    this.props.patchMatch(this.props.matchId, {
      action: 'map-pick',
      map: e.target.dataset.id,
    }, this.refresh);
  };

  render() {
    const {
      gameMaps, bracket, bracketRound, bracketMaps, match, matchLeadership,
      matchMaps,
    } = this.props;

    if (!gameMaps || !bracket || !bracketRound || !matchLeadership) return null;

    const myTeamId = (matchLeadership.size === 1) ?
      matchLeadership.first() : null;

    const picked = [];
    const pickedIds = [];
    for (const mm of matchMaps.toJS()) {
      picked.push(
        <div key={`p${mm.id}`} className={this.cn({
          d: 'avail',
          m: [
            'picked',
            mm.teamId && (mm.teamId === match.get('teamX') ? 'x' : 'y'),
            mm.isBan ? 'ban' : 'pick',
          ],
        })}>
          <div className={this.cn({ d: 'mapOuter' })}>
            <div className={this.cn({
              d: 'map',
              m: mm.isBan && 'nope',
            })}>
              {gameMaps.get(mm.gameMapId).get('name')}
            </div>
          </div>
        </div>
      );
      pickedIds.push(mm.gameMapId);
    }

    let rawProcedure = bracket.get('mapVetoProcedure');
    if (bracketRound.get('mapVetoProcedure') !== '') {
      rawProcedure = bracketRound.get('mapVetoProcedure');
    }

    const procedure = rawProcedure.split(' ').reduce(
      (r, x) => (x === '' ? r : r.concat(x)), [],
    );

    let avail;
    const areMapsReady = match.get('areMapsReady');
    if (
      !areMapsReady &&
      procedure.length &&
      picked.length < procedure.length
    ) {
      const action = procedure[picked.length];
      let letter = action[0];
      const isBan = letter.toUpperCase() !== letter;
      letter = letter.toUpperCase();
      const [isX, isY] = [letter === 'X', letter === 'Y'];
      const canPick = (
        (isX && myTeamId === match.get('teamX')) ||
        (isY && myTeamId === match.get('teamY'))
      );
      const sub = +action[1];
      const hasSub = !isNaN(sub);
      avail = (
        <div className={this.cn({
          d: 'avail',
          m: [
            isBan ? 'ban' : 'pick',
            isX && 'x', isY && 'y',
          ],
        })}>{
          bracketMaps && bracketMaps
          .filter(x => !pickedIds.includes(x.get('gameMapId')))
          .sortBy((x) => {
            const gm = gameMaps.get(x.get('gameMapId'));
            return (gm && gm.get('name')) || x.get('gameMapId');
          })
          .toJS().map((x) => {
            const gm = gameMaps.get(x.gameMapId);
            const name = (gm && gm.get('name')) || x.gameMapId;
            const isPickable = !hasSub || x.subPool === sub;
            return (
              <div key={`a${x.id}`} className={this.cn({
                d: 'mapOuter',
                m: [isX && 'x', isY && 'y'],
              })}>
                <button data-id={x.gameMapId} className={this.cn({
                  d: 'map',
                  m: isPickable ? (canPick && 'pointer') : 'nope',
                })} onClick={canPick && isPickable && this.pick}>
                  {name} <span
                    className={this.cn({ d: 'mapPool' })}
                    title="map sub-pool"
                  >
                    {x.subPool}
                  </span>
                </button>
              </div>
            );
          })
        }</div>
      );
    }

    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ d: 'inner' })}>
          {picked}
          {!!picked.length && avail &&
            <div className={this.cn({ d: 'rule' })} />}
          {avail}
          {rawProcedure !== '' && <div className={this.cn({ d: 'rule' })} />}
          {rawProcedure !== '' && <div className={this.cn({ d: 'procedure' })}>
            <div className={this.cn({ d: 'procedureTitle' })}>
              THE PROCEDURE
            </div>
            <div
              className={this.cn({ d: 'procedureList' })}
            >{procedure.map((action, i) => {
              let letter = action[0];
              const isBan = letter.toUpperCase() !== letter;
              letter = letter.toUpperCase();
              const [isX, isY, isR] = [
                letter === 'X',
                letter === 'Y',
                letter === 'R',
              ];
              const sub = +action[1];
              const hasSub = !isNaN(sub);
              const isDone = areMapsReady || i < picked.length;
              return (
                <div key={`proc${i}`}><div className={this.cn({
                  d: 'procedureItem',
                  m: [isBan ? 'ban' : 'pick', isDone && 'done'],
                })}>
                  {isR && 'auto-'}{isBan ? 'ban' : 'pick'} {hasSub && <span
                    className={this.cn({ d: 'mapPool' })}
                    title="map sub-pool"
                  >{sub}</span>}
                  {!isR && <div className={this.cn({
                    d: 'procedureWho',
                    m: [isX && 'x', isY && 'y', isDone && 'done'],
                  })} />}
                </div></div>
              );
            })}</div>
          </div>}
        </div>
      </div>
    );
  }
}
