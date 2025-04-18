import { useParams } from 'react-router-dom';
import {
  TaskRunDetail,
  TabOption,
} from '../../workflow-runs-v1/$run/v2components/step-run-detail/step-run-detail';
import invariant from 'tiny-invariant';

export default function TaskRun() {
  const params = useParams();

  invariant(params.run);

  return (
    <div className="px-8">
      <TaskRunDetail taskRunId={params.run} defaultOpenTab={TabOption.Output} />
    </div>
  );
}
