import { setState, getState } from '../../../state/state.js';
import { fetchProfile } from '../../profile/fetchProfile.js';
import { toggleActionRequest } from '../../profile/api.js';
import Notify from '../../../components/ui/Notify.js';

interface ToggleLabels { on: string; off: string }

interface ToggleActionOptions {
  entityId: string | number;
  entityType?: string;
  button: HTMLButtonElement | null;
  apiPath: string;
  labels?: ToggleLabels;
  actionName?: string;
}

export async function toggleAction({ entityId, entityType = 'user', button, apiPath, labels = { on: 'Active', off: 'Inactive' }, actionName = 'action' }: ToggleActionOptions): Promise<void> {
  if (!getState('token')) {
    Notify('Please log in first.', { type: 'warning', duration: 3000, dismissible: true });
    return;
  }

  if (!button) {
    Notify('Action button not found.', { type: 'info', duration: 3000, dismissible: true });
    return;
  }

  const isActive = button.dataset.active === 'true';
  const httpMethod = isActive ? 'DELETE' : 'PUT';
  const apiEndpoint = `${apiPath}${entityId}`;

  const originalText = button.textContent;
  const wasActive = isActive;

  button.disabled = true;
  button.textContent = isActive ? labels.off : labels.on;
  button.dataset.active = String(!isActive);

  try {
    const response = (await toggleActionRequest(apiEndpoint, httpMethod)) as Response | undefined;

    if (response && typeof response.ok === 'boolean' && !response.ok) {
      throw new Error(`Server responded with status ${response.status}`);
    }

    button.disabled = false;

    if (entityType === 'user') {
      const updatedProfile = await fetchProfile();
      if (updatedProfile) {
        setState({ userProfile: updatedProfile }, true);
      }
    }

    const actionText = !wasActive ? actionName : `un${actionName}`;
    Notify(`You have ${actionText} this ${entityType}.`, { type: 'success', duration: 3000, dismissible: true });
  } catch (error) {
    button.textContent = originalText;
    button.dataset.active = String(wasActive);
    button.disabled = false;

    const err = error as Error;
    console.error(`Error toggling ${actionName}:`, error);
    Notify(`Failed to update ${actionName}: ${err.message || 'Unknown error'}`, { type: 'error', duration: 3000, dismissible: true });
  }
}

export function toggleFollow(userId: string | number, followButton: HTMLButtonElement | null): Promise<void> {
  return toggleAction({
    entityId: userId,
    entityType: 'user',
    button: followButton,
    apiPath: '/follow/',
    labels: { on: 'Unfollow', off: 'Follow' },
    actionName: 'followed'
  });
}
