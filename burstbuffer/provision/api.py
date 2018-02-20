# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.

import os

from burstbuffer.registry import api

ASSIGNED_SLICES_KEY = "bufferhosts/assigned_slices/%s"
ALL_SLICES_KEY = "bufferhosts/all_slices/%s/%s"

FAKE_DEVICE_COUNT = 12
FAKE_DEVICE_ADDRESS = "nvme%sn1"
FAKE_DEVICE_SIZE_BYTES = int(1.5 * 2 ** 40)  # 1.5 TB


def _get_local_hardware():
    fake_devices = []
    for i in range(FAKE_DEVICE_COUNT):
        fake_devices.append(FAKE_DEVICE_ADDRESS % i)
    return fake_devices


def _update_data(data):
    for key, value in data.iter():
        api._etcdctl("put '%s' '%s'" % (key, value))


def _refresh_slices(hostname, hardware):
    slices_info = {}
    for device in hardware:
        key = ALL_SLICES_KEY % (hostname, device)
        slices[key] = FAKE_DEVICE_SIZE_BYTES
    _update_data(slices_info)


def _get_assigned_slices(hostname):
    prefix = ASSIGNED_SLICES_KEY % hostname
    raw_assignments = api._get_all_with_prefix(prefix)
    current_devices = _get_local_hardware()

    assignemnts = {}
    for key in raw_assignments:
        device = key[(len(prefix) + 1):]
        if device not in current_devices:
            raise Exception("assignment to unknown device %s!!" % device)
        assignemnts[device] = raw_assignments[key]
    return assignemnts


def startup(hostname):
    all_slices = _get_local_hardware()
    _refresh_slices(hostname, all_slices)

    return _get_assigned_slices(hostname)


def _get_env():
    return os.environ


def _get_event_info():
    env = _get_env()

    event_type = env["ETCD_WATCH_EVENT_TYPE"]
    revision = env["ETCD_WATCH_REVISION"]
    key = env["ETCD_WATCH_KEY"]
    value = env['ETCD_WATCH_VALUE']

    return dict(
        event_type=event_type,
        revision=revision,
        key=key,
        value=value)


def event(hostname):
    event_info = _get_event_info()
    print(event_info)
    return _get_assigned_slices(hostname)
