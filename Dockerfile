# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause
FROM scratch
COPY vra-cli /
ENTRYPOINT [ "/vra-cli" ]